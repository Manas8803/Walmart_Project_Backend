const AWS = require("aws-sdk");

exports.handler = async function (event, context) {
  let connections;
  const body = JSON.parse(event.body);
  const discountOffer = body.data.discount_offer;
  connections = body.data.connection_ids;

  if (!discountOffer || !connections || connections.length === 0) {
    return {
      statusCode: 400,
      message: "Missing discount_offer or connection_ids in the request body.",
    };
  }

  const callbackAPI = new AWS.ApiGatewayManagementApi({
    apiVersion: "2018-11-29",
    endpoint: event.requestContext.domainName + "/" + event.requestContext.stage,
  });

  let errors = [];

  for (let conn of connections) {
    const connectionId = conn;
    if (connectionId !== event.requestContext.connectionId) {
      try {
        await callbackAPI
          .postToConnection({ 
            ConnectionId: connectionId, 
            Data: JSON.stringify({ discount_offer: discountOffer }) 
          })
          .promise();
      } catch (e) {
        console.log(`Error posting to connection ${connectionId}:`, e);
        errors.push(connectionId);
      }
    }
  }

  if (errors.length > 0) {
    return {
      statusCode: 500,
      message: `Failed to post to some connections: ${errors.join(", ")}`,
    };
  }

  return { statusCode: 200 };
};