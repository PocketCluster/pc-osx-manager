# README

This is a self contained package and it should be as it links to "C" and
being compiled to a static binary. Do not import this package from others.

### GO buildmode

Go 1.7.5 has following buildmode

```
-buildmode=pie [Not supported in amd64]
    Build the listed main packages and everything they import into
    position independent executables (PIE). Packages not named
    main are ignored.

-buildmode=exe
    Build the listed main packages and everything they import into
    executables. Packages not named main are ignored.

-buildmode=default
    Listed main packages are built into executables and listed
    non-main packages are built into .a files (the default
    behavior).

-buildmode=archive
    Build the listed non-main packages into .a files. Packages named
    main are ignored.

-buildmode=shared
    Combine all the listed non-main packages into a single shared
    library that will be used when building with the -linkshared
    option. Packages named main are ignored.
```

- - -

```
-buildmode=c-archive
    Build the listed main package, plus all packages it imports,
    into a C archive file. The only callable symbols will be those
    functions exported using a cgo //export comment. Requires
    exactly one main package to be listed.

-buildmode=c-shared
    Build the listed main packages, plus all packages that they
    import, into C shared libraries. The only callable symbols will
    be those functions exported using a cgo //export comment.
    Non-main packages are ignored.
```