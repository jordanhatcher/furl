## Furl

Furl is a project to bundle a target executable along with the shared libraries it depends on, as a single portable executable. 
The target executable and libraries are compressed and embedded within a wrapper executable. The wrapper decompresses and saves 
the embedded files to a temporary directory, and uses the embedded dynamic linker to run the embedded executable.

### Building

```
make build executable=<executable> linker=<linker> objectfiles-dir=<objectfile directory>
```