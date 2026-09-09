# threepo

<img src = https://raw.githubusercontent.com/FrankKair/threepojs/master/assets/threepo-logo.png width="25%" height="25%"/>

Transform a localization spreadsheet (XLSX or CSV) into JSON.

Expects a header row with a `key` column and one column per locale:

| key      | en       | pt         | sv    |
|----------|----------|------------|-------|
| hello    | hello    | olá        | hej   |
| computer | computer | computador | dator |

Output:

```json
{
  "hello":      {"en": "hello",    "pt": "olá",        "sv": "hej"},
  "computador": {"en": "computer", "pt": "computador", "sv": "dator"}
}
```

Dot-separated keys produce nested objects: `home.title` becomes `{ "home": { "title": { ... } } }`.

## Install

```
go install github.com/FrankKair/threepo@latest
```

Or build from source:

```
make build
```

## Usage

```
threepo <file.xlsx|file.csv>
```

```
threepo strings.xlsx
threepo strings.csv
threepo strings.xlsx | jq .key   # pipe to jq
threepo strings.xlsx > out.json  # write to file
```
