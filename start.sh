 #!/usr/bin/env bash

amixer scontrols 'Speaker' 100% unmute
amixer scontrols 'Mic' 100% unmute
amixer scontrols 'Auto Gain Control' 100% unmute

exec ./main
