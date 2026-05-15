// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Package builder transforms profile content during the prebuild pipeline.
//
// A builder implements the [Builder] interface, whose [Builder.Apply] method
// receives an [Option] describing the profile being processed and the current
// profile content as a string, and returns the transformed content. Builders
// are stateless and may be invoked concurrently across profiles.
//
// Builders are composed into a pipeline via [Builders], created with
// [NewRunner]. The pipeline exposes a fluent [Builders.Add] for registration
// and [Builders.Run] to apply every registered builder, in order, to a
// profile:
//
//	r := builder.NewRunner(cfg).
//		Add(builder.NewABI3()).
//		Add(builder.NewComplain()).
//		Add(builder.NewUserspace())
//
//	out, err := r.Run(file, content)
//
// Order matters: each builder sees the output of the previous one, so a
// builder that rewrites attachments must run before one that depends on the
// resolved form, and mode-changing builders (complain, enforce) should run
// after content-shaping ones.
package builder
