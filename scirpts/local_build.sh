PROJECT_NAME="ece428_mp2"
BRANCH="master"
TARGET_FILE_NAME="gossip"
TARGET_PATH="./cmd"
FOLDER_NAME=${PROJECT_NAME}
OUTPUT="execution_out"
EXECUTE_NAME="gossip_service"

cd ~/go

# make sure the target project folder exists
if [ ! -d "$FOLDER_NAME" ]; then
  git clone -b ${BRANCH} git@gitlab.engr.illinois.edu:dl58/${PROJECT_NAME}.git
fi

# Build the executable file
cd ./${FOLDER_NAME} && \
git pull && \
go build -o ${EXECUTE_NAME} ${TARGET_PATH}/${TARGET_FILE_NAME}.go  && \
echo "source file built successfully!"

# Make sure there are no older versions of programs running in the VM, and then execute the target program in the background.
PID=`pgrep ${EXECUTE_NAME}`
if [ "${PID}" != "" ]; then
  pkill ${EXECUTE_NAME}
fi
touch ${OUTPUT} && \
nohup ./${EXECUTE_NAME} > ./${OUTPUT} 2>&1 &

exit 0
