# launch
launch is a small Go tool that allows you to run VS Code's configurations (described in launch.json) from command line and without launching VS Code.

Currently only launching executable specified by **program** with environment variables specified in **environment** is supported.

## Usage
#### Examine the output
`launch <configuration name>`

#### Execute by piping to sh compatible shell:
`launch <configuration name> | sh`

