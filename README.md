## Introduction

A simple and naive tool for importing Trello tasks to GitHub issues. This tool was designed for a specific use case, hence make sure to read through the assumptions before running the application. Feel free to make a fork of the project and modify it to suit your needs.

## Assumptions

- Tasks are imported to a GitHub project where no issues have been created before.
- Issues are closed if they are in list "Done" or "Testing" in Trello.
- User is assignee for all issues.
- User is the owner of the repository.
- All attachments from Trello are images.
- The tool is run without interruptions.
- The directory that contains the executable, and from where the tool is executed, contains a well formatted config.json file.
- Trello board does not contain archived elements (this has not been tested).

## Usage

Create a personal access token in GitHub, and export Trello board as JSON.

Create file config.json in the following format:

```
{
    "gitHubToken": "<OAuth token>",
    "trelloBoardJsonPath": "<exported json file>",
    "user": "<GitHub user>",
    "repository": "<GitHub repository>"
}
```

Fill in required data and place the file in the same directory where the executable is located.

Run the created executable.

## Features

The following task data is imported:
- Task title.
- Description, contains the following data:
    - Task description
    - Checklists
    - Actions
    - Comments
    - Images (only links)
    - Link to task in Trello
- Closes the issue if it is in list "Testing" or "Done" in Trello.
