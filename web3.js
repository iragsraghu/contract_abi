function web3go() {
	// functions for web3.js
	var Contract = require('web3-eth-contract');
	Contract.setProvider('https://mainnet.infura.io/v3/7ba7186d11d24eddbf53996feb6dbabf');
	var getJSON = require('get-json');
	var web3 = require('web3');


	address = "0xEDd27C961CE6f79afC16Fd287d934eE31a90D7D1"
	var url = "https://api.etherscan.io/api?module=contract&action=getabi&address=0xEDd27C961CE6f79afC16Fd287d934eE31a90D7D1&apikey=M2QBXVREAF3JQH23I2G58FAQSF7189Q3XF"

	var request = require('request')

	request({url: url, json: true}, function (error, response, body) {
		var buf = Buffer.from(JSON.stringify(body.result));
		var stringData = JSON.parse(buf.toString());
		var abi = JSON.parse(stringData);
		var contract = new Contract(abi, address);
		var bigNumberInputAmount = web3.utils.toWei('0.1', 'ether');
		inputdata = contract.methods.enter(bigNumberInputAmount.toString()).encodeABI();
		console.log(inputdata);
	})
}

web3go();