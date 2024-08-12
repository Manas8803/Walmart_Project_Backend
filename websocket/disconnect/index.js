const AWS = require("aws-sdk");
const ddb = new AWS.DynamoDB.DocumentClient();

exports.handler = async function (event, context) {
  const connectionId = event.requestContext.connectionId;

  try {
    // Scan the user table to find the item with the given connectionId
    const scanParams = {
      TableName: process.env.USER_TABLE_ARN,
      FilterExpression: "contains(connection_id, :connectionId)",
      ExpressionAttributeValues: {
        ":connectionId": connectionId
      }
    };

    const scanResult = await ddb.scan(scanParams).promise();

    if (scanResult.Items.length === 0) {
      return {
        statusCode: 404,
        body: JSON.stringify({
          message: "No user found with the given connectionId",
        }),
      };
    }

    const userToUpdate = scanResult.Items[0];
    
    // Remove the current connectionId from the connection_id array
    const updatedConnectionIds = userToUpdate.connection_id.filter(
      (id) => id !== connectionId
    );

    // Update the user item with the new connectionIds array
    await ddb.update({
      TableName: process.env.USER_TABLE_ARN,
      Key: { user_id: userToUpdate.user_id },
      UpdateExpression: "SET connection_id = :updatedConnectionIds",
      ExpressionAttributeValues: {
        ":updatedConnectionIds": updatedConnectionIds,
      },
    }).promise();

    return {
      statusCode: 200,
      body: JSON.stringify({ message: "Disconnected successfully" }),
    };
  } catch (err) {
    console.error("Error disconnecting:", err);
    return {
      statusCode: 500,
      body: JSON.stringify({ error: "Error disconnecting" }),
    };
  }
};