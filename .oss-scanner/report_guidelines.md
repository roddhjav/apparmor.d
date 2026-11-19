# How we would like reports written

## Layout and length

- Under one page, with these sections in this order: Summary, Chain, Rules,
  Impact, Suggested fix.
- One escape chain per report. If several profiles share the root cause (an
  abstraction, a tunable), report the shared file once and list the affected
  profiles.

## What every report should state

- The starting profile: which confined program the attacker controls, and why
  it is plausible that it gets compromised (network facing, parses user files,
  runs as root...).
- The chain, step by step: each step names the profile in effect and the
  exact rule that allows it, with its file and line in `apparmor.d/`
  (source, not `.build/`), including rules that come from abstractions.
- The end position: unconfined, which profile, which uid, which data.
- The distributions concerned when rules depend on `#aa:only`,
  `#aa:exclude` or on the ABI (4 or 5).

## Style

- Plain text, short sentences. No speculation presented as fact: say
  explicitly when a step is assumed and not shown from the rules.
- Title in the form `<profile or abstraction>: <what escapes>`.
- Do not use the em dash character.
