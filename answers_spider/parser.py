import sys
import time

for line in sys.stdin:
  try:
    sys.stdout.write(line)
    sys.stderr.write(line)
    sys.stdout.flush()
  except:
    pass
  time.sleep(10)

