#!/usr/bin/env python3
# SPDX-License-Identifier: GPL-2.0-only

# KNOWN ISSUES:
# No guards for file type - expects AppArmor
# Diffirent suggestions for single line are mutually exclusive
# Suggestion could point to changed profile name, based on other suggestion

import sys
import argparse
import pathlib
import shlex
import json
from copy import deepcopy

try:
    from apparmor.regex         import *
    from apparmor.aa            import is_skippable_file
    from apparmor.rule.file     import FileRule, FileRuleset
    from apparmor.common        import convert_regexp
    try:
        from apparmor.rule.variable import separate_vars
    except ImportError:
        from apparmor.aa            import separate_vars

    LIBAPPARMOR = True

except ImportError:
    LIBAPPARMOR = False

def sanitizeProfileName(name):

    if name.startswith('/') or name.startswith('@{'):
        name = pathlib.Path(name).stem

    if ' ' in name:
        name = re.sub(r'\s+', '-', name)

    return name

def makeLocalIdentity(nestingStacker_):

    newStacker = []
    for i in nestingStacker_:
        i = sanitizeProfileName(i)
        newStacker.append(i)

    identity = '_'.join(newStacker)  # separate each (sub)profile identity with underscores

    return identity

def getCurrentProfile(stacker):

    if stacker:
        profile = stacker[-1]
    else:
        profile = None

    return profile

def handleFileMessages(l, file, profile, lineNum):

    wholeFileAccessProfiles = (
#        '',
    )
    suggestOwner = (  # TODO: switch to AARE
        r'^@{HOME}/',
        r'^/home/\w+/',
        r'^@{run}/user/@{uid}/',
        r'^/run/user/\d+/',
        r'^@{tmp}/',
        r'^/tmp/',
        r'^/var/tmp/',
        r'^/dev/shm/',
    )

    lG = l.groupdict()
    reason_ = None
    if lG.get('path'):
        if lG.get('path').startswith('/**') and profile not in wholeFileAccessProfiles:  # false positives
            severity_   = 'ERROR'
            reason_     = 'Whole filesystem access is too broad'
            suggestion_ = None

        for r in suggestOwner:
            if re.match(r, lG.get('path')) and not lG.get('owner'):
                indentRe = re.match(r'^\s+', l.group())
                if indentRe:
                    indent = indentRe.group()
                else:
                    indent = ''

                severity_   = 'NOTICE'
                reason_     = "'owner' is likely required"
                suggestion_ = indent + 'owner ' + l.group().lstrip()
                break

    elif lG.get('bare_file') and profile not in wholeFileAccessProfiles:
        severity_   = 'ERROR'
        reason_     = 'Whole filesystem access is too broad'
        suggestion_ = None

    if reason_:  # something matched
        msg = {'filename':   file,
               'profile':    profile,
               'severity':   severity_,
               'line':       lineNum,
               'reason':     reason_,
               'suggestion': suggestion_}
    else:
        msg = None

    return msg

def readApparmorFile(fullpath):
    '''AA file could contain multiple AA profiles'''
    headers = (
        '# apparmor.d - Full set of apparmor profiles',
        '# Copyright (C) ',
        '# SPDX-License-Identifier: GPL-2.0-only',
    )

    file_data = {}
    fileVars = {}
    nestingStacker = []
    duplicateProfilesCounter = []
    localExists = {}
    localExists_eol = {}
    messages = []
    exceptionMsg = None
    line = None
    gotAbi = False
    gotHeaders = {}
    gotAttach = False
    isAfterProfileStart = False
    lastLineNum = None
    try:
        with open(fullpath, 'r') as f:
            for n,line in enumerate(f, start=1):
                if isAfterProfileStart:
                    isAfterProfileStart = False
                    expectedIndent = len(nestingStacker) * '  '
                    indentRe = re.match(r'^\s+', line)
                    if indentRe:
                        indent = indentRe.group()
                    else:
                        indent = ''

                    if indent != expectedIndent:
                        spacesCount  = len(nestingStacker) * 2
                        nestingCount = len(nestingStacker)
                        messages.append({'filename':   fullpath,
                                         'profile':    getCurrentProfile(nestingStacker),
                                         'severity':   'WARNING',
                                         'line':       n,
                                         'reason':     f"Expected {spacesCount} spaces for {nestingCount} nesting",
                                         'suggestion': f"{expectedIndent}{line.lstrip()}"})

                if line.endswith(' \n'):
                    messages.append({'filename':   fullpath,
                                     'profile':    getCurrentProfile(nestingStacker),
                                     'severity':   'WARNING',
                                     'line':       n,
                                     'reason':     "Redundant trailing whitespace",
                                     'suggestion': line.rstrip()})

                if '\t' in line:
                    messages.append({'filename':   fullpath,
                                     'profile':    getCurrentProfile(nestingStacker),
                                     'severity':   'WARNING',
                                     'line':       n,
                                     'reason':     "Tabs are not allowed",
                                     'suggestion': line.replace('\t', '  ')})

                if len(gotHeaders) < 3 and not nestingStacker:
                    for nH,i in enumerate(headers):
                        if line.startswith(i):
                            gotHeaders[nH] = True

                if   RE_ABI.search(line):
                    gotAbi = line

                elif RE_PROFILE_START.search(line) or RE_PROFILE_HAT_DEF.search(line):
                    isAfterProfileStart = True
                    m = parse_profile_start_line(line, fullpath)
                    if m.get('profile'):
                        nestingStacker.append(m.get('profile'))  # set early

                    if m.get('attachment') == '@{exec_path}' and not gotAttach:  # can be only singular
                        gotAttach = True
                        
                    profileMsg = {'filename':   fullpath,
                                  'profile':    getCurrentProfile(nestingStacker),
                                  'severity':   'WARNING',
                                  'line':       n,
                                  'reason':     "A short named profile must be defined",
                                  'suggestion': None}
                    if   m.get('plainprofile'):
                        messages.append(profileMsg)
                    elif m.get('namedprofile'):
                        if m.get('namedprofile').startswith('/'):
                            messages.append(profileMsg)

                    if m.get('flags'):
                        m['flags'] = set(shlex.split(m.pop('flags').replace(',', '')))
                        if 'complain' in m['flags']:
                            messages.append({'filename':   fullpath,
                                             'profile':    getCurrentProfile(nestingStacker),
                                             'severity':   'WARNING',
                                             'line':       n,
                                             'reason':     "'complain' flag must be defined in 'dists/flags'",
                                             'suggestion': None})
                    else:
                        m['flags'] = set()

                    if m.get('profile'):
                        duplicateProfilesCounter.append(m.get('profile'))
                        profileIdentity = '//'.join(nestingStacker)
                        file_data[profileIdentity] = m

                elif RE_PROFILE_VARIABLE.search(line):
                    lineV = RE_PROFILE_VARIABLE.search(line).groups()
 
                    name = strip_quotes(lineV[0])
                    operation = lineV[1]
                    val = separate_vars(lineV[2])
                    if fileVars.get(name):
                        fileVars[name].update(set(val))
                        if operation == '=':
                            messages.append({'filename':   fullpath,
                                             'profile':    getCurrentProfile(nestingStacker),
                                             'severity':   'DEGRADED',
                                             'line':       n,
                                             'reason':     "Tunable must be appended with '+='",
                                             'suggestion': None})
                    else:
                        fileVars[name] = set(val)
                        if operation == '+=':
                            messages.append({'filename':   fullpath,
                                             'profile':    getCurrentProfile(nestingStacker),
                                             'severity':   'DEGRADED',
                                             'line':       n,
                                             'reason':     "Tunable must be defined with '='",
                                             'suggestion': None})

                elif RE_INCLUDE.search(line):
                    if nestingStacker:
                        profileIdentity = '//'.join(nestingStacker)
                        localIdentity = makeLocalIdentity(nestingStacker)
                        localValue = f'include if exists <local/{localIdentity}>'  # commented out will also match
                        if localValue in line:
                            localExists[profileIdentity] = localValue

                # Handle file entries
                elif RE_PROFILE_FILE_ENTRY.search(line):
                    lineF = RE_PROFILE_FILE_ENTRY.search(line)
                    fileMsg = handleFileMessages(lineF, fullpath, getCurrentProfile(nestingStacker), n)
                    if fileMsg:
                        messages.append(fileMsg)

                elif RE_PROFILE_END.search(line):
                    if getCurrentProfile(nestingStacker):
                        if not nestingStacker:
                            messages.append({'filename':   fullpath,
                                             'profile':    None,
                                             'severity':   'DEGRADED',
                                             'line':       n,
                                             'reason':     "Unbalanced parenthesis?",  # not fully covered
                                             'suggestion': None})
                        else:
                            profileIdentity = '//'.join(nestingStacker)
                            localExists_eol[profileIdentity] = n
                            del nestingStacker[-1]  # remove last

                lastLineNum = n

    except PermissionError:
        exceptionMsg = 'Unable to read the file (PermissionError)'

    except UnicodeDecodeError:
        exceptionMsg = 'Unable to read the file (UnicodeDecodeError)'

    except FileNotFoundError:
        exceptionMsg = 'No such file or directory (FileNotFoundError)'

    if exceptionMsg:
        messages.append({'filename':   fullpath,
                         'profile':    None,
                         'severity':   'NOTICE',
                         'line':       None,
                         'reason':     exceptionMsg,
                         'suggestion': None})

    # Ensure proper header is present
    if len(gotHeaders) < 3:
        combinedHeader = '\n'.join(headers)
        messages.append({'filename':   fullpath,
                         'profile':    None,
                         'severity':   'WARNING',
                         'line':       1,
                         'reason':     'No proper header',
                         'suggestion': combinedHeader})

    # Ensure ABI is present
    changeAbi = False
    abi = 'abi <abi/4.0>,'
    if gotAbi:
        if gotAbi.strip() != abi:
            changeAbi = True
    else:
        changeAbi = True

    if changeAbi:
        messages.append({'filename':   fullpath,
                         'profile':    None,
                         'severity':   'WARNING',
                         'line':       None,
                         'reason':     'ABI is required',
                         'suggestion': abi})

    # Ensure singular '@{exec_path}'
    if not gotAttach:
        messages.append({'filename':   fullpath,
                         'profile':    None,
                         'severity':   'WARNING',
                         'line':       None,
                         'reason':     "'@{exec_path}' must be defined as main path attachment",
                         'suggestion': None})

    # Ensure trailing vim syntax
    if line:
        trailingSyntax = '# vim:syntax=apparmor\n'
        if line != trailingSyntax:
            messages.append({'filename':   fullpath,
                             'profile':    None,
                             'severity':   'WARNING',
                             'line':       lastLineNum,
                             'reason':     'No trailing syntax hint',
                             'suggestion': trailingSyntax})

    # Assign variables to profile attachments as paths and assign filenames
    for p,d in deepcopy(file_data).items():
        file_data[p]['filename'] = fullpath
        attachment = d.get('attachment')
        if attachment:
            if attachment.startswith('@{'):
                if fileVars.get(attachment):
                    file_data[p]['attach_paths'] = fileVars[attachment]  # incoming set
                else:
                    messages.append({'filename':   fullpath,
                                     'profile':    p,
                                     'severity':   'ERROR',
                                     'line':       None,
                                     'reason':     f"Unknown global variable as profile attachment: {attachment}",
                                     'suggestion': None})

            else:
                if isinstance(file_data[p].get('attachment'), set):
                    raise ValueError("Expecting 'str' or 'None', not 'set'")
                file_data[p]['attach_paths'] = {file_data[p]['attachment']}

    # Check if profile block does not have corresponding 'local' include
    for p,d in file_data.items():
        if not localExists.get(p):  # not found previously
            if '//' in p:
                identity = p.split('//')
            else:
                identity = [p]

            localIdentity = makeLocalIdentity(identity)
            filename = file_data[p]['filename']
            messages.append({'filename':   filename,
                             'profile':    p,
                             'severity':   'WARNING',
                             'line':       localExists_eol.get(p),  # None? Unbalanced parenthesis?
                             'reason':     "The (sub)profile block does not have expected 'local' include",
                             'suggestion': f'include if exists <local/{localIdentity}>'})

    # Track multiple definitions inside single file
    for profile in duplicateProfilesCounter:
        counter = duplicateProfilesCounter.count(profile)
        if counter >= 2:
            messages.append({'filename':   fullpath,
                             'profile':    profile,
                             'severity':   'DEGRADED',
                             'line':       None,
                             'reason':     "Profile has been defined {counter} times in the same file",
                             'suggestion': None})

    return (messages, file_data)

def findAllProfileFilenames(profile_dir):

    profiles = set()
    for path in pathlib.Path(profile_dir).iterdir():
        if path.is_file() and not is_skippable_file(path):
            profiles.add(path.resolve())

    # Not default, dig deeper
    if not profiles:
        nestedDirs = (
            'groups',
            'profiles-a-f',
            'profiles-g-l',
            'profiles-m-r',
            'profiles-s-z',
        )
        for d in nestedDirs:
            dirpath = pathlib.Path(pathlib.Path(profile_dir).resolve(), pathlib.Path(d))
            for p in dirpath.rglob("*"):
                if p.is_file():
                    profiles.add(p)

    return profiles

def handleArgs():
    """DEGRADED are purposed for fatal errors - when the profile set will fail to load entirely"""

    allSeverities = ['DEBUG', 'NOTICE', 'WARNING', 'ERROR', 'CRITICAL', 'DEGRADED']
    aaRoot = '/etc/apparmor.d'

    parser = argparse.ArgumentParser()
    parser.add_argument('-d', '--aa-root-dir', action='store',
                        default=aaRoot,
                        help='Target different AppArmor root directory rather than default')
    parser.add_argument('-p', '--profile', action='append',
                        help='Handle only specified profile')
#    parser.add_argument('-s', '--severity', action='append',
#                        choices=allSeverities,
#                        help='Handle only specified severity event')

    args = parser.parse_args()

#    if not args.severity:
#        args.severity = allSeverities

    return args

def main(argv):

    args = handleArgs()

    messages = []

    profile_dir = args.aa_root_dir
    if not args.profile:
        profiles = findAllProfileFilenames(profile_dir)
    else:
        profiles = set()
        for p in args.profile:
            absolutePath = pathlib.Path(p).resolve()
            profiles.add(absolutePath)

    profile_data = {}
    for path in sorted(profiles):
        if not is_skippable_file(path):
            readApparmorFile_Out = readApparmorFile(path)
            profilesInFile = readApparmorFile_Out[1]
            messages.extend(readApparmorFile_Out[0])
            profile_data.update(profilesInFile)

    for m in messages:
        if m.get('suggestion'):
            if m['suggestion'].endswith('\n'):
                m['suggestion'] = m.get('suggestion').removesuffix('\n')
        m['filename'] = str(m.get('filename'))
        print(json.dumps(m, indent=2))

    if messages:
        sys.exit(1)

    return None

if __name__ == '__main__':

    if not LIBAPPARMOR:
        raise ImportError(f"""Can't find 'python3-apparmor' package! Install with:
$ sudo apt install python3-apparmor""")

    main(sys.argv)
