var stuff = {
	authToken: null
}

$(document).ready(function() {
	req("help", null, handleHelp)
})

function handleHelp(res) {
	console.log(res)
	var $cmdlist = $('#cmdlist')
	$cmdlist.html('')
	for (x in res) {
		var cmd = res[x];
		var dt = $('<dt>'+x+'</dt>')
		var dd = $('<dd>'+cmd.Desc+'<br></dd>')
		for (arg in cmd.Args) {
			var desc = cmd.Args[arg]
			var s = '<span class="cmdlabel">'+
				nl2br(escapeHtml(desc))+'</span><br>'

			if (arg.indexOf('__') !== 0) {
				s += '<input name="'+arg+'" type="text">'
			}
			s += '<br>'
			dd.append(s)
		}
		var el = $('<form class="cmdform">')
		el.data('cmd', x)
		el.append(dt, dd, '<input type="submit">')
		$cmdlist.append(el)
	}
	$('form.cmdform').submit(function() {
		args = {}
		$(this).find('input').each(function() {
			if (this.type === "submit") return
			if (this.value === null || this.value === "") return
			args[this.name] = this.value
		})
		console.log(args)
		var cmd = $(this).data('cmd')
		req(cmd, args, handleRes.bind(null,cmd))
		return false
	})
}

function handleRes(cmd, res) {
	if (cmd === "login") {
		handleLogin(res)
	} else if (cmd === "help") {
		handleHelp(res)
	}
}

function handleLogin(res) {
	if (res.status !== "OK") {
		console.error(res)
		return
	}
	stuff.authToken = res.auth_token
	req("help", null, handleHelp)
}

function req(cmd, args, callback) {
	if (stuff.authToken) {
		args = args || {}
		args['auth_token'] = stuff.authToken
	}
	$.ajax({
		url: apiurl+cmd,
		data: args,
		success: callback,
		mimeType: "application/json",
		error: function(err) {
			console.error(err);
		},
		type: "GET"
	})
}


var entityMap = {
	"&": "&amp;",
	"<": "&lt;",
	">": "&gt;",
	'"': '&quot;',
	"'": '&#39;',
	"/": '&#x2F;'
};

function escapeHtml(string) {
	return String(string).replace(/[&<>"'\/]/g, function (s) {
		return entityMap[s];
	});
}
function nl2br(string) {
	return string.replace(/\n/g, "<br>")
}
