import os
from dotenv import load_dotenv

def mongoDumpAllData(dumpDirName, mongoDumpServerIP, mongoDumpDBName):
    cliMongodump = "mongodump" \
        + " -o " + dumpDirName \
        + " --db " + mongoDumpDBName \
        + " -h " + mongoDumpServerIP

    os.system(cliMongodump)

def mongoDump(dumpDirName, mongoDumpServerIP, mongoDumpDBName, mongoDumpCollections):
    for i in range(len(mongoDumpCollections)):
        cliMongodump = "mongodump" \
            + " --collection  " + mongoDumpCollections[i] \
            + " -o " + dumpDirName \
            + " --db " + mongoDumpDBName \
            + " -h " + mongoDumpServerIP

        os.system(cliMongodump)

def mongoRestore(dumpDirName, mongoDumpDBName, mongoRestoreServerIP, mongoRestoreDBName):
    cliMongoRestore = 'mongorestore' \
        + ' '+ dumpDirName +'/' + mongoDumpDBName \
        + ' --db ' + mongoRestoreDBName \
        + ' -h ' + mongoRestoreServerIP \
        + ' --drop'

    os.system(cliMongoRestore)

if __name__ == '__main__':
    load_dotenv()

    dumpDirName = os.getenv("MONGO_DUMP_DIR_NAME")

    isMongoDumpAllData = os.getenv("MONGO_DUMP_ALL_DATA")

    mongoDumpServerIP = os.getenv("MONGO_DUMP_SERVER_IP")
    mongoDumpDBName = os.getenv("MONGO_DUMP_DB_NAME")
    mongoDumpCollections = os.getenv("MONGO_DUMP_COLLECTION").split(',')

    mongoRestoreDBName = os.getenv("MONGO_RESTORE_DB_NAME")
    mongoRestoreServerIP = os.getenv("MONGO_RESTORE_SERVER_IP")

    if isMongoDumpAllData == "true":
        mongoDumpAllData(dumpDirName, mongoDumpServerIP, mongoDumpDBName)
    else:
        mongoDump(dumpDirName, mongoDumpServerIP, mongoDumpDBName, mongoDumpCollections)

    print("")
    mongoRestore(dumpDirName, mongoDumpDBName, mongoRestoreServerIP, mongoRestoreDBName)

    os.system("rm -rf " + dumpDirName)