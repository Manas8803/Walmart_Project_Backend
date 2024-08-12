const AWS = require("aws-sdk");
const ddb = new AWS.DynamoDB.DocumentClient();

exports.handler = async function (event, context) {
    const user_id = event.queryStringParameters?.user_id;

    if (!user_id) {
        return {
            statusCode: 200,
            body: JSON.stringify({ message: "Connected successfully" }),
        };
    }

    try {
        // Retrieve the existing user data
        const userData = await ddb.get({
            TableName: process.env.USER_TABLE_ARN,
            Key: { user_id: user_id },
        }).promise();

        if (!userData.Item) {
            return {
                statusCode: 404,
                body: JSON.stringify({
                    message: `User with id '${user_id}' not found`,
                }),
            };
        }

        // Get the existing connection_ids or initialize an empty array
        let connectionIds = userData.Item.connection_ids || [];

        // Add the new connectionId to the array
        connectionIds.push(event.requestContext.connectionId);

        // Update the user item with the new connectionIds array
        await ddb.update({
            TableName: process.env.USER_TABLE_ARN,
            Key: { user_id: user_id },
            UpdateExpression: "SET connection_ids = :connectionIds",
            ExpressionAttributeValues: {
                ":connectionIds": connectionIds,
            },
        }).promise();

        return {
            statusCode: 200,
            body: JSON.stringify({ message: "Connected successfully" }),
        };
    } catch (err) {
        console.error("Error connecting:", err);
        return {
            statusCode: 500,
            body: JSON.stringify({ error: "Error connecting" }),
        };
    }
};