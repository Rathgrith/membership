MACHINE_LIST=(
                "fa23-cs425-4801.cs.illinois.edu"
                "fa23-cs425-4802.cs.illinois.edu"
                "fa23-cs425-4803.cs.illinois.edu"
#               "fa23-cs425-4804.cs.illinois.edu"
#                "fa23-cs425-4805.cs.illinois.edu"
#                "fa23-cs425-4806.cs.illinois.edu"
#                "fa23-cs425-4807.cs.illinois.edu"
#                "fa23-cs425-4808.cs.illinois.edu"
#                "fa23-cs425-4809.cs.illinois.edu"
#                "fa23-cs425-4810.cs.illinois.edu"
               )

USER_NAME=""
PASSWORD=""
while getopts "u:p:" opt; do
  case $opt in
   u)
     USER_NAME=$OPTARG
    ;;
   p)
     PASSWORD=$OPTARG
     ;;
   *)
     echo "invalid arg:$OPTARG" && exit
     ;;
  esac
done

for machine_host in ${MACHINE_LIST[*]}
  do
      ssh ${USER_NAME}@${machine_host} "bash -s" < kill.sh
  done
