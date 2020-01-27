# Stackdriver Logging Configuration #
___

## Useful tools: ##

 1. [Regexr](https://regexr.com/) for validation checks / writing regexp
 2. [Fluentular](https://fluentular.herokuapp.com/) for checking fluentd statements with sample log entries
 3. Ruby [strftime](https://apidock.com/ruby/Time/strftime) reference for formatting timestamp fields
 
 ___
 
 * To add/modify a log file configuration for stackdriver, see the `configs` directory. Do not directly edit the `.conf` files found in the `logging-config` directory.
 
 * Once changes have been made, run `main.go` to generate the configuration files before pushing to GitLab.
 
 * As a reference guide, see the `examples` directory for templates that will cover most use cases
 
 * All ` \ ` characters in the configuration files require escaping (e.g. ` \\ `), see existing logs for examples
 
 * Optional features follow indenting (e.g. the 'modifyFields' option is only available under the 'recordTransformer' heading)
 ___
 
 ## Breakdown of config file ##
 
 ```
 # All placeholders are marked in the format: {PLACEHOLDER} 
 # This config WILL NOT run as it includes every possible option available

 # All fields are required unless specified
 
 logType:
   source:          "tail"                                                           # Currently only the 'tail' input plugin is supported
   system:          "{SYSTEM THE LOG FILE IS FROM}"                                  # OPTIONAL - Quick identifier for the log
   path:            "{PATH TO THE LOCATION OF THE LOG FILE ON THE HOST MACHINE}"     # Location of the file being read by stackdriver
   posPath:         "{PATH TO THE LOCATION OF THE POS FILE ON THE HOST MACHINE}"     # Pos file created by stackdriver
   tag:             "{INDENTIFICATION TAG FOR GCP LOG VIEWER}"                       # Identifier for this log in GCP e.g. 'example.audit.log'
   timeFormat:      "{RUBY STRFTIME TIME FORMAT}"                                    # e.g. %H:%M:%S for a log starting '11:59:59'
   parseType:       "{FLUENTD PARSE TYPE}"                                           # either 'regexp', 'multiline' or 'multi_format' 
   formatFirstLine: "{REGEX TO TRACK START OF NEW ENTRY (USUALLY THE TIME/DATE)}"    # MULTILINE parseType ONLY - regex that will determine the start of a log entry
                                                                                     #
 recordTransformer:                                                                  # OPTIONAL - allows for transformation of the log before sending it to GCP
   removeKeys:                                                                       # OPTIONAL - removes specified fields from the log
     - "{FIELD TO REMOVE}"                                                           # Field name to be removed
   modifyFields:                                                                     # OPTIONAL - allows for ruby string manipulation on field content
     - field: "{FIELD TO CHANGE}"                                                    # Field to modify (after parsing)
       modify: "{RUBY STRING TRANSFORMATION SNIPPET}"                                # Ruby snippet to modify string content 
                                                                                     #
 fields:                                                                             # REGEXP/MULTILINE parseType ONLY
   - name: "{FIRST FIELD IN LOG}"                                                    # Name of the field - may be GCP reserved (see below)
     regex: "{REGEX TO CAPTURE FIRST FIELD OF LOG}"                                  # Regular expression to match the content of this field
     delimiter: "{REGEX TO DENOTE END OF FIRST FIELD}"                               # Regular expression to match the delimter between two fields
   - name: "{NEXT FIELD IN LOG}"                                                     #    
     regex: "{REGEX TO CAPTURE CURRENT FIELD}"                                       #
     delimiter: "{REGEX TO DENOTE END OF CURRENT FIELD}"                             #
                                                                                     #
 multiFormat:                                                                        # MULTI_FORMAT parseType ONLY
   - type: "regexp"                                                                  # ParseType to use for current format
   - fields:                                                                         # Name of field (uses same structure as REGEXP/MULTILINE formats above)
     - name: "{FIRST FIELD IN LOG}"                                                  #
     - regex: "{REGEX TO CAPTURE FIRST FIELD OF LOG}"                                #
     - delimiter: "{REGEX TO DENOTE END OF FIRST FIELD}"                             #
 ```

___

## Parse types and use cases ##

1. `REGEXP` - the simplest parse type. This format reads each new line of the log file as a new entry, and uses regular expressions to format the raw text in each entry. Use this for simple single-line, single-format logs.

2. `MULTILINE` - behaves much in the same way as the `REGEXP` parser, with the additional functionality of being able to handle log entries that span multiple lines in the log file. If the log outputs things like stacktraces, use this parser.

3. `MULTI_FORMAT` - is designed to handle logs that do not share the same format but are compiled in the same log file. This parser allows you to describe multiple parsers to match the various formats present in the log file. It will use the first parser that matches the log entry it is currently trying to read. 

___

## Main Reserved fields ##

Reserved fields are special field names that GCP will look for to match certain properties of the log entry. For example, any field named `severity` will be used to set the severity indicator (e.g. DEBUG, INFO, ERROR etc.) that is displayed alongside log entries in the GCP log viewer.

| __Field Name__ | __Maps To / sets__                           |
| :-------------:|:--------------------------------------------:|
| time           | `created_time`                               |
| severity       | `severity indicator`                         |
| message        | `the main message displayed in the log`      |

Please see [here](https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry) for more details. 

___

## Pos files ##

`.pos` files are created by the stackdriver agent when it first opens a log file for reading. A `.pos` file acts like a bookmark; the logging agent stores the line in the log file that it read most recently, so that in the event of a crash or a restart, it doesn't start from the beginning of the log file which would waste time and may result in duplicate entries.

Aside from telling the logging agent where to create the pos file, there is nothing else we need to do. The agent is able to handle log rotations and recover from crashes providing the `.pos` file is present.

A `.pos` file won't be created if the log file doesn't exist or is empty. 
