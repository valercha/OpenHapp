'use strict';
'require form';
'require view';

return view.extend({
	render: function() {
		var m, s, o;

		m = new form.Map(
			'openhapp',
			_('OpenHapp Configuration'),
			_('Configure the OpenHapp runtime daemon.')
		);

		s = m.section(
			form.NamedSection,
			'main',
			'openhapp',
			_('Runtime')
		);

		o = s.option(
			form.Flag,
			'enabled',
			_('Enabled'),
			_('Enable the OpenHapp runtime daemon.')
		);
		o.rmempty = false;
		o.default = '1';

		o = s.option(
			form.ListValue,
			'engine',
			_('Engine'),
			_('Select the proxy engine used by OpenHapp.')
		);
		o.value('sing-box', _('Sing-box'));
		o.rmempty = false;
		o.default = 'sing-box';

		o = s.option(
			form.ListValue,
			'mode',
			_('Mode'),
			_('Runtime routing mode.')
		);
		o.value('proxy', _('Proxy'));
		o.rmempty = false;
		o.default = 'proxy';

		o = s.option(
			form.ListValue,
			'log_level',
			_('Log level'),
			_('Logging verbosity for the OpenHapp daemon.')
		);
		o.value('debug', _('Debug'));
		o.value('info', _('Info'));
		o.value('warn', _('Warning'));
		o.value('error', _('Error'));
		o.rmempty = false;
		o.default = 'info';

		o = s.option(
			form.Value,
			'listen',
			_('Listen address'),
			_('Local control/listen address used by OpenHapp.')
		);
		o.rmempty = false;
		o.default = '127.0.0.1:0';
		o.validate = function(section_id, value) {
			if (!value)
				return _('Listen address must not be empty.');

			var parts = value.split(':');

			if (parts.length !== 2 || !parts[0] || !parts[1])
				return _('Use host:port format, for example 127.0.0.1:0.');

			var port = Number(parts[1]);

			if (!Number.isInteger(port) || port < 0 || port > 65535)
				return _('Port must be between 0 and 65535.');

			return true;
		};

		o = s.option(
			form.Flag,
			'autostart',
			_('Autostart'),
			_('Start OpenHapp automatically during system boot.')
		);
		o.rmempty = false;
		o.default = '1';

		o = s.option(
			form.Value,
			'subscription',
			_('Subscription'),
			_('VPN subscription URL. Leave empty when not configured.')
		);
		o.rmempty = true;

		return m.render();
	}
});
