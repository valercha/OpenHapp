'use strict';
'require view';
'require rpc';

var callStatus = rpc.declare({
	object: 'openhapp',
	method: 'status'
});

var callEngineInfo = rpc.declare({
	object: 'openhapp',
	method: 'engine_info'
});

var callStart = rpc.declare({
	object: 'openhapp',
	method: 'start'
});

var callStop = rpc.declare({
	object: 'openhapp',
	method: 'stop'
});

return view.extend({
	title: _('OpenHapp Actions'),

	load: function () {
		return Promise.all([
			callStatus().catch(function () { return {}; }),
			callEngineInfo().catch(function () { return {}; })
		]);
	},

	render: function (data) {
		var status = data[0] || {};
		var engineInfo = data[1] || {};
		var busy = false;

		var runtimeState = status.running ? _('running') : _('stopped');
		var engineName = engineInfo.name || '-';
		var engineState = engineInfo.running ? _('running') : _('stopped');
		var engineAvailability = engineInfo.available ? _('available') : _('unavailable');

		function runAction(action) {
			if (busy)
				return;

			busy = true;

			action()
				.then(function () {
					return window.location.reload();
				})
				.catch(function (err) {
					window.alert(
						err && err.message
							? err.message
							: _('Action failed')
					);
				})
				.finally(function () {
					busy = false;
				});
		}

		return E('div', { 'class': 'cbi-map' }, [
			E('h2', {}, _('OpenHapp Actions')),

			E('div', { 'class': 'cbi-section' }, [
				E('p', {}, [
					E('strong', {}, _('OpenHapp daemon: ')),
					runtimeState
				]),
				E('p', {}, [
					E('strong', {}, _('Engine: ')),
					engineName
				]),
				E('p', {}, [
					E('strong', {}, _('Engine state: ')),
					engineState
				]),
				E('p', {}, [
					E('strong', {}, _('Engine availability: ')),
					engineAvailability
				])
			]),

			E('div', { 'class': 'cbi-section' }, [
				E('button', {
					'class': 'btn cbi-button cbi-button-action',
					'click': function () {
						runAction(callStart);
					}
				}, _('Start OpenHapp')),

				' ',

				E('button', {
					'class': 'btn cbi-button cbi-button-reset',
					'click': function () {
						runAction(callStop);
					}
				}, _('Stop OpenHapp'))
			])
		]);
	}
});
