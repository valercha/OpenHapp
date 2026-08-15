'use strict';
'require view';
'require rpc';

var callStatus = rpc.declare({
	object: 'openhapp',
	method: 'status'
});

var callConfig = rpc.declare({
	object: 'openhapp',
	method: 'config'
});

var callManifest = rpc.declare({
	object: 'openhapp',
	method: 'manifest'
});

var callEngineInfo = rpc.declare({
	object: 'openhapp',
	method: 'engine_info'
});

return view.extend({
	title: _('OpenHapp Dashboard'),

	load: function () {
		return Promise.all([
			callConfig().catch(function () { return {}; }),
			callStatus().catch(function () { return {}; }),
			callManifest().catch(function () { return {}; }),
			callEngineInfo().catch(function () { return {}; })
		]);
	},

	render: function (data) {
		var config = data[0] || {};
		var status = data[1] || {};
		var manifest = data[2] || {};
		var engineInfo = data[3] || {};

		var enabled = config.enabled ? '1' : '0';
		var engine = config.engine || 'xray';
		var mode = config.mode || 'proxy';
		var listen = config.listen || '127.0.0.1:0';
		var subscription = config.subscription || '';

		var runtimeState = status.running ? _('running') : _('stopped');
		var manifestVersion = manifest.version || '-';
		var manifestUpdatedAt = manifest.updated_at || '-';
		var manifestName = manifest.name || 'OpenHapp';
		var engineVersion = engineInfo.version || '-';
		var engineBinary = engineInfo.binary || '-';
		var engineConfig = engineInfo.config || '-';
		var engineWorkdir = engineInfo.workdir || '-';
		var engineAvailable = engineInfo.available ? _('available') : _('unavailable');
		var engineRunning = engineInfo.running ? _('running') : _('stopped');

		return E('div', { 'class': 'cbi-map' }, [
			E('h2', {}, _('OpenHapp Dashboard')),
			E('div', { 'class': 'cbi-section' }, [
				E('p', {}, [ E('strong', {}, _('Runtime: ')), enabled === '1' ? _('enabled') : _('disabled') ]),
				E('p', {}, [ E('strong', {}, _('Daemon state: ')), runtimeState ]),
				E('p', {}, [ E('strong', {}, _('Engine: ')), engine ]),
				E('p', {}, [ E('strong', {}, _('Engine version: ')), engineVersion ]),
				E('p', {}, [ E('strong', {}, _('Engine state: ')), engineRunning ]),
				E('p', {}, [ E('strong', {}, _('Engine availability: ')), engineAvailable ]),
				E('p', {}, [ E('strong', {}, _('Engine binary: ')), engineBinary ]),
				E('p', {}, [ E('strong', {}, _('Engine config: ')), engineConfig ]),
				E('p', {}, [ E('strong', {}, _('Engine workdir: ')), engineWorkdir ]),
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
