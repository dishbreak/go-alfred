# Alfred Workflows in Go

This is a work-in-progress, but it aims to implement the JSON interface that Alfred for macOS expects when it executes commands in its workflows.

The Version 2 rewrite attempts to make the library much more lightweight. It leaves items like configuration and execution control to the caller and focuses on sending appropriate responses back to Alfred using its JSON format.

Version 2 requires Alfred 5 or later in order to take advantage of things like Alfred's built-in caching options.
