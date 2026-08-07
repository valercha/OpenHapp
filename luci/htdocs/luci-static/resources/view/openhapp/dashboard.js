'use strict';
'require view';
'require uci';
'require rpc';

var callStatus = rpc.declare({
	object: 'openhapp',
	method: 'status'
});

var callManifest = rpc.declare({
	object: 'openhapp',
	method: 'manifest'
});

return view.extend({
	title: _('OpenHapp Dashboard'),

	load: function () {
		return Promise.all([
			uci.load('openhapp'),
			callStatus().catch(function () { return {}; }),
			callManifest().catch(function () { return {}; })
		]);
	},

	render: function (data) {
		var uciData = data[0] || {};
		var status = data[1] || {};
		var manifest = data[2] || {};
		var section = uciData.main || null;
		var enabled = section && section.enabled ? section.enabled : '0';
		var engine = section && section.engine ? section.engine : 'xray';
		var mode = section && section.mode ? section.mode : 'proxy';
		var listen = section && section.listen ? section.listen : '127.0.0.1:0';
		var subscription = section && section.subscription ? section.subscription : '';

		var runtimeState = status.running ? _('running') : _('stopped');
		var manifestVersion = manifest.version || '-';
		var manifestUpdatedAt = manifest.updated_at || '-';
		var manifestName = manifest.name || 'OpenHapp';

		return E('div', { 'class': 'cbi-map' }, [
			E('h2', {}, _('OpenHapp Dashboard')),
			E('div', { 'class': 'cbi-section' }, [
				E('p', {}, [ E('strong', {}, _('Runtime: ')), enabled === '1' ? _('enabled') : _('disabled') ]),
				E('p', {}, [ E('strong', {}, _('Daemon state: ')), runtimeState ]),
				E('p', {}, [ E('strong', {}, _('Engine: ')), engine ]),
				E('p', {}, [ E('strong', {}, _('Mode: ')), mode ]),
				E('p', {}, [ E('strong', {}, _('Listen: ')), listen ]),
				E('p', {}, [ E('strong', {}, _('Subscription: ')), subscription || _('none') ]),
				E('p', {}, [ E('strong', {}, _('Manifest: ')), manifestName ]),
				E('p', {}, [ E('strong', {}, _('Manifest version: ')), manifestVersion ]),
				E('p', {}, [ E('strong', {}, _('Manifest updated: ')), manifestUpdatedAt ])
			])
		]);
	}
});