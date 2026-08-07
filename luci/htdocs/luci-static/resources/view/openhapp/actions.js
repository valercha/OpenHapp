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

		return E('div', { 'class': 'cbi-map' }, [
			E('h2', {}, _('OpenHapp Actions')),
			E('div', { 'class': 'cbi-section' }, [
				E('p', {}, [ E('strong', {}, _('Daemon state: ')), runtimeState ])
			]),
			E('div', { 'class': 'cbi-section' }, [
				E('button', {
					'class': 'btn cbi-button cbi-button-action',
					'click': function () {
						return callStart().then(function () { window.location.reload(); });
					}
				}, _('Start')),
				' ',
				E('button', {
					'class': 'btn cbi-button cbi-button-reset',
					'click': function () {
						return callStop().then(function () { window.location.reload(); });
					}
				}, _('Stop'))
			])
		]);
	}
});
