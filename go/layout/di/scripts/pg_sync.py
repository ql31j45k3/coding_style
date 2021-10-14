import os
from dotenv import load_dotenv

def pgDumpAllData(dumpDirName, pgDumpServerIP, pgDumpPort, pgDumpUserName, pgDumpDBName):
    cliDump = "pg_dump" \
        + " -h " + pgDumpServerIP \
        + " -p " + pgDumpPort \
        + " -U " + pgDumpUserName \
        + " -Fc " + pgDumpDBName \
        + " > " + dumpDirName + ".dump"

    print(os.popen(cliDump).readlines())

def pgDumpOnlySchema(dumpDirName, pgDumpServerIP, pgDumpPort, pgDumpUserName, pgDumpDBName):
    cliDump = "pg_dump" \
        + " -h " + pgDumpServerIP \
        + " -p " + pgDumpPort \
        + " -U " + pgDumpUserName \
        + " --schema-only" \
        + " -Fc " + pgDumpDBName \
        + " > " + dumpDirName + ".dump"

    print(os.popen(cliDump).readlines())

def pgDump(dumpDirName, pgDumpServerIP, pgDumpPort, pgDumpUserName, pgDumpDBName, pgDumpTables):
    for i in range(len(pgDumpTables)):
        cliDump = "pg_dump" \
            + " -h " + pgDumpServerIP \
            + " -p " + pgDumpPort \
            + " -U " + pgDumpUserName \
            + " -t " + pgDumpTables[i]\
            + " -Fc " + pgDumpDBName \
            + " > " + dumpDirName + "_" + pgDumpTables[i] + ".dump"

        print(os.popen(cliDump).readlines())

def pgRestore(dumpDirName, pgDumpServerIP, pgDumpPort, pgDumpUserName, pgRestoreDBName):
    cliRestore = "pg_restore" \
        + " -h " + pgDumpServerIP \
        + " -p " + pgDumpPort \
        + " -U " + pgDumpUserName \
        + " -c" \
        + " -d " + pgRestoreDBName + " " + dumpDirName + ".dump"

    print(os.popen(cliRestore).readlines())

def pgRestoreTables(dumpDirName, pgDumpServerIP, pgDumpPort, pgDumpUserName, pgRestoreDBName, pgDumpTables):
    for i in range(len(pgDumpTables)):
        cliRestore = "pg_restore" \
            + " -h " + pgDumpServerIP \
            + " -p " + pgDumpPort \
            + " -U " + pgDumpUserName \
            + " -c" \
            + " -d " + pgRestoreDBName + " " + dumpDirName + "_" + pgDumpTables[i] + ".dump"

        print(os.popen(cliRestore).readlines())

if __name__ == '__main__':
    load_dotenv()

    dumpDirName = os.getenv("PG_DUMP_DIR_NAME")
    pgDumpServerIP = os.getenv("PG_DUMP_SERVER_IP")
    pgDumpPort = os.getenv("PG_DUMP_PORT")
    pgDumpUserName = os.getenv("PG_DUMP_USERNAME")

    pgDumpDBName = os.getenv("PG_DUMP_DB_NAME")

    pgRestoreDBName = os.getenv("PG_RESTORE_DB_NAME")

    pgDumpTables = os.getenv("PG_DUMP_TABLE").split(',')

    isPGDumpAllData = os.getenv("PG_DUMP_ALL_DATA")
    isPGDumpOnlySchema = os.getenv("PG_DUMP_ONLY_SCHEMA")

    if isPGDumpAllData == "true":
        pgDumpAllData(dumpDirName, pgDumpServerIP, pgDumpPort, pgDumpUserName, pgDumpDBName)
        print("")
        pgRestore(dumpDirName, pgDumpServerIP, pgDumpPort, pgDumpUserName, pgRestoreDBName)
    elif isPGDumpOnlySchema == "true":
        pgDumpOnlySchema(dumpDirName, pgDumpServerIP, pgDumpPort, pgDumpUserName, pgDumpDBName)
        print("")
        pgRestore(dumpDirName, pgDumpServerIP, pgDumpPort, pgDumpUserName, pgRestoreDBName)
    else:
        pgDump(dumpDirName, pgDumpServerIP, pgDumpPort, pgDumpUserName, pgDumpDBName, pgDumpTables)
        pgRestoreTables(dumpDirName, pgDumpServerIP, pgDumpPort, pgDumpUserName, pgRestoreDBName, pgDumpTables)

    os.system("rm -rf " + dumpDirName + "*.dump")