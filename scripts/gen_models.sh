#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
(
  cd "$DIR/.."
  go install github.com/volatiletech/sqlboiler/v4@latest
  sqlboiler --no-rows-affected --no-auto-timestamps --no-hooks --no-tests --no-context --add-global-variants --add-panic-variants --wipe mysql
)
