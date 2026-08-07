'use strict';
'require view';
'require rpc';

var callStatus = rpc.declare({
	object: 'openhapp',
	method: 'status'
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
		return callStatus().catch(function () { return {}; });
	},

	render: function (data) {
		var status = data || {};
		var runtimeState = status.running ? _('running') : _('stopped');
		var busy = false;

		function runAction(action) {
			if (busy)
				return;

			busy = true;
			action()
				.then(function () {
					window.location.reload();
				})
				.catch(function (err) {
					window.alert(err && err.message ? err.message : _('Action failed'));
				})
				.finally(function () {
					busy = false;
				});
		}

		return E('div', { 'class': 'cbi-map' }, [
			E('h2', {}, _('OpenHapp Actions')),
			E('div', { 'class': 'cbi-section' }, [
				E('p', {}, [ E('strong', {}, _('Daemon state: ')), runtimeState ])
			]),
			E('div', { 'class': 'cbi-section' }, [
				E('button', {
					'class': 'btn cbi-button cbi-button-action',
					'click': function () { runAction(callStart); }
				}, _('Start')),
				' ',
				E('button', {
					'class': 'btn cbi-button cbi-button-reset',
					'click': function () { runAction(callStop); }
				}, _('Stop'))
			])
		]);
	}
});