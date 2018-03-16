import {
  OrganizationClient
} from '../set-up/client';
import config from '../set-up/config';
const orgId = process.env.ORGID || "org.FitCoinOrg";
<<<<<<< HEAD:containers/blockchain/fitcoinBackend/utils/peer.js
//const clientArray = config.peers.filter(obj => obj.peer.org === orgId).map(obj => new OrganizationClient(config.channelName, config.orderer, obj.peer, obj.ca, obj.admin));
//const numberOfClients = process.env.WORKERCLIENTS || 1;
const workerClients = config.peers.filter(obj => obj.peer.org === orgId).map(obj => new OrganizationClient(config.channelName, config.orderer, obj.peer, obj.ca, obj.admin));
=======
const workerClients = config.peers.filter(obj => obj.peer.org === orgId).map(obj => new OrganizationClient(config.channelName, config.orderer, obj.peer, obj.ca, obj.admin));
var eventEmitterClient = null;
var clientArray = {
  workers: workerClients,
};
if(process.env.EVENTEMITTER === "true") {
  eventEmitterClient = new OrganizationClient(workerClients[0]._channelName, workerClients[0]._ordererConfig, workerClients[0]._peerConfig, workerClients[0]._caConfig, workerClients[0]._admin);
  clientArray.eventEmitter = eventEmitterClient;
}
>>>>>>> fdce2c1ef2bac2d962da8c44e0b25ef0ccf231ef:containers/blockchain/backend/utils/peer.js
export async function initiateClient() {
  try {
    for(var i = 0; i < workerClients.length; i++) {
      await workerClients[i].login();
    }
    await Promise.all([workerClients.map(obj => obj.initEventHubs())]);
    if(process.env.EVENTEMITTER === "true") {
      await eventEmitterClient.login();
      await eventEmitterClient.initEventHubs();
    }
  } catch(e) {
    console.log('Fatal error logging into blockchain organization clients!');
    console.log(e);
    process.exit(-1);
  }
}
exports.clients = clientArray;