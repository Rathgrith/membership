EXECUTE_NAME="gossip_service"

PID=`pgrep ${EXECUTE_NAME}`
if [ "${PID}" != "" ]; then
  pkill ${EXECUTE_NAME}
fi