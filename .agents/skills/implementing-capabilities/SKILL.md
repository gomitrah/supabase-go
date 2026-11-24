---
name: implementing-capabilities
description: Guides what the SDK does against a Supabase service, as distinct from the Go shape of its API - whether a feature or option belongs in the SDK at all, what a request sends, what the client validates, which defaults it applies and how it decodes responses, plus the capability declarations in sdk-compliance.yaml. Use when adding or changing a feature, an option, request or response handling or a capability declaration, when reviewing such a change or when asked about the Supabase SDK engineering principles.
---

The following documents will be helpful:

| Document | For |
| -------- | --- |
| [Supabase SDK engineering principles](https://github.com/supabase/sdk/blob/main/docs/principles.md) | The principles every official Supabase client SDK follows: thin clients, logic and validation in the service, service-owned defaults and decoding that accepts new values, plus the registry and spec-first ways of working. |
| [`decisions.md`](../../../decisions.md) | The records "`sdk-compliance.yaml` is sparse" and "Compliance capability ids never appear in doc comments". |

Read the principles before writing or reviewing such a change, and cite a principle by a link to its heading.
The conventions in this repository win for local matters such as commands, layout and style.

The Go shape of the surface a capability adds is the `designing-public-interfaces` skill's territory.
