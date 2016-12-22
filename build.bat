#!/bin/sh

gox -os "darwin linux windows" -arch "amd64" -output "{{.Dir}}_{{.OS}}_{{.Arch}}/{{.Dir}}"

for %%i in (darwin linux windows) DO (
    MKDIR samaritan_%%i_amd64\web\dist
    MKDIR samaritan_%%i_amd64\custom
    COPY web\dist\ samaritan_%%i_amd64\web\dist\
    COPY config.ini samaritan_%%i_amd64\custom\config.ini
    COPY README.md samaritan_%%i_amd64\README.md
    COPY LICENSE samaritan_%%i_amd64\LICENSE
)
