import os
import argparse
from pickle import TRUE
import subprocess
import zipfile

ROOTDIR = "~/Development/AWS/todo"

def buildBootstrap(dir):
    os.chdir(dir)

    build_command = "GOOS=linux GOARCH=amd64 go build -o bootstrap bootstrap.go"
    subprocess.run(build_command, shell=True)

    zipBootstrap()

def zipBootstrap():
    fileName = "bootstrap"
    zipName  = "function.zip"
    if not os.path.exists(fileName):
        print(f"Error: Failed to compile bootstrap.")
        return


    with zipfile.ZipFile(zipName, "w") as zippy:
        zippy.write(fileName)

def buildSAM():
    build_command = "sam build"
    subprocess.run(build_command, shell=True)

def deploySAM(mode):
    if mode == "GUIDED":
        build_command = "deploy --guided"
    else:
        build_command = "deploy"

    subprocess.run(build_command, shell=True)

def runSAMLocally():
    build_command = "sam local start-api"
    subprocess.run(build_command, shell=True)

def buildLambdaFunctions():
    current_dir = os.getcwd()
    todoDirs = [
        "createTodo",
        "deleteTodo",
        "getTodo",
        "updateTodo",
    ]
    for child_dir in todoDirs:
        if os.path.exists(child_dir) and os.path.isdir(child_dir):
            print(f" --> Building {child_dir}'s bootstrap.go file..")
            buildBootstrap(child_dir)
            print(f" --> {child_dir}'s Bootstrap Complete.")
        else:
            print(f"Child Directory ./\"{child_dir}\" doesn't exist")
        os.chdir(current_dir)

def main(args):
    subprocess.run("clear;", shell=True)
    if args.subcommand == "lam" and args.sam:
        print("Lambda -> sam")
        buildLambdaFunctions()
        buildSAM()
        return

    if args.subcommand == "lam":
        print("Lambda")
        buildLambdaFunctions()
        return

    if args.subcommand == "sam" and args.lam:
        print("Sam -> Lambda")
        buildSAM()
        buildLambdaFunctions()
        return

    if args.subcommand == "sam" and not args.lam:
        print("Sam")
        buildSAM()
        return

    if args.subcommand == "all" and args.local:
        print("Lambda -> Sam -> Local")
        buildLambdaFunctions()
        buildSAM()
        runSAMLocally()
        return

    if args.subcommand == "all" and args.deploy_guided:
        print("Lambda -> Sam -> Deploy --guided")
        buildLambdaFunctions()
        buildSAM()
        deploySAM(mode="GUIDED")
        return

    if args.subcommand == "all" and args.deploy:
        print("Lambda -> Sam -> Deploy")
        buildLambdaFunctions()
        buildSAM()
        deploySAM(mode="")
        return

    if args.subcommand == "all":
        print(" --> Missing required Subcommand")
        print("    -> --local")
        print("    -> --deploy")
        print("    -> --deploy-guided")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Build script for building both our Lambda binaries, and our AWS Servees"
    )
    subparsers = parser.add_subparsers(dest="subcommand")

    buildType = subparsers.add_parser("lam", help="Builds our Lambda Functions only")
    buildType.add_argument(
        "--sam",
        help="Builds SAM after we build our Lambda Functions. Does not run our AWS service",
        action="store_true"
    )

    buildType = subparsers.add_parser("sam", help="Builds SAM without running or Deploying")
    buildType.add_argument(
        "--lam",
        help="Build our Lambda Functions after we build our AWS Service. Does not run our AWS Service",
        action="store_true"
    )

    parser_all = subparsers.add_parser(
        "all",
        help="Run all all build scripts, including running either locally or in Deployment"
    )
    parser_all.add_argument(
        "--local",
        action="store_true",
        help="After build is complete, we run our AWS Service Locally."
    )
    parser_all.add_argument(
        "--deploy",
        action="store_true",
        help="After build is complete, we deploy our Service to AWS."
    )
    parser_all.add_argument(
        "--deploy-guided",
        action="store_true",
        help="After build is complete, we deploy our Service, with Guidance to AWS."
    )

    args = parser.parse_args()
    main(args)
