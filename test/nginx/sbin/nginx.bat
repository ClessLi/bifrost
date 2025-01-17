@ECHO OFF
IF "%~1"=="-t" (
  ECHO "check failure" 1>&2
  exit 1
)
ECHO "pass"
exit 0