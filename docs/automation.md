# Automated Video Creation with n8n

This document outlines how to use the n8n workflow to automate the video creation pipeline in the ASSOS project.

## Accessing the n8n Dashboard

The n8n service is included in the `docker-compose.yml` file and runs alongside the other services. To access the n8n dashboard:

1.  Start the application using `docker-compose up`.
2.  Open your web browser and navigate to `http://localhost:5678`.

You will be prompted to set up an admin user on your first visit.

## Importing the Workflow

The core video creation workflow is defined in the `n8n/workflow.json` file. To import it into n8n:

1.  From the n8n dashboard, click on "Workflows" in the left-hand sidebar.
2.  Click the "Import from File" button.
3.  Select the `n8n/workflow.json` file from the project root.

The "Automated Video Creation Pipeline" workflow will be imported and ready to use.

## Understanding the Workflow

The workflow consists of the following steps:

1.  **Manual Trigger:** The workflow is started manually by clicking the "Execute Workflow" button. This is for testing and demonstration purposes. In a production environment, this could be replaced with a schedule, a webhook, or another trigger.

2.  **Create Video Record:** This step makes an HTTP POST request to the `api-gateway` to create a new video record in the database. This action also triggers the `ai-service` to begin generating a script for the video.
    *   **Note:** This step requires an authentication token. You will need to create a "Header Auth" credential in n8n with the name `assos-api-auth` and provide a valid JWT for an authenticated user.

3.  **Wait for AI Script:** This is a simple "Wait" node that pauses the workflow for one minute. This is a temporary solution to allow the `ai-service` enough time to generate the script. In a more advanced setup, this could be replaced with a webhook system for instant continuation.

4.  **Trigger Video Processing:** After the wait, this step makes another HTTP POST request to the `api-gateway` to begin the video rendering process. It uses the `video_id` returned from the "Create Video Record" step to specify which video to process.

## Running the Workflow

To run the workflow:

1.  Ensure you have created the `assos-api-auth` credential in n8n.
2.  Open the "Automated Video Creation Pipeline" workflow.
3.  Click the "Execute Workflow" button in the top-right corner.

You can monitor the progress of the workflow in the n8n dashboard and check the application logs to see the backend services being triggered. The final rendered video will be available at the URL stored in the `videos` table in the database.