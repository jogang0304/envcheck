# envcheck

Package which exports function "envcheck.Load()". This function loads variables from dotenv file into os.env and checks that all variables pass to config written in .env.yaml.

.env.yaml structure:
```yaml
vars:
  - name: string
    required: bool
    type: enum{"string", "int", "float", "bool", "any"}
    default_value: value of specified type
    pattern: regex string (if type is string)
  - name: string
    ...
```