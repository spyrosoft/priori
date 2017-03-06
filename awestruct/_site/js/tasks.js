var tasks = [];

$('.new-task input.task').select();
populate_tasks();

$(window).on('hashchange', hash_changed);


function hash_changed() {
	populate_tasks();
	var parent_id = get_parent_id_from_url_hash();
	$('form.new-task .parent-id').first().val(parent_id);
}


function populate_tasks() {
	$('.tasks').html('');
	var post_data = {'action': 'tasks'};
	var parent_id = get_parent_id_from_url_hash();
	if (typeof parent_id !== 'undefined') {
		post_data['parent-id'] = parent_id;
	}
	api_call(
		post_data,
		'populate-tasks'
	);
}

success_response_callbacks['populate-tasks'] = function(post_data, tasks_data) {
	if (tasks_data['parent'] !== undefined) {
		$('#page-title').html(Belt.escapeHTML(tasks_data['parent']));
	}
	if (tasks_data['parent-id'] !== undefined) {
		$('a.parent-task').attr('href', '/#' + tasks_data['parent-id']);
	} else {
		$('a.parent-task').attr('href', '/#');
	}
	if (tasks_data['tasks'] === undefined || tasks_data['tasks'] === '') {return;}
	
	try {
		tasks = JSON.parse(tasks_data['tasks']);
	} catch(e) {
		display_error('An error occurred while retrieving your tasks. An admin has been notified.');
		notify_admin(e);
		return;
	}
	display_tasks();
};

function display_tasks() {
	$('.tasks').html('');
	tasks.sort(sort_tasks);
	for (var i in tasks) {
		append_task(
			tasks[i]['task'],
			tasks[i]['id'],
			tasks[i]['short_term'],
			tasks[i]['long_term'],
			tasks[i]['urgency'],
			tasks[i]['difficulty']
		);
	}
}

function sort_tasks(task_1, task_2) {
	var task_1_total
		= parseInt(task_1['short_term'])
		+ parseInt(task_1['long_term'])
		+ parseInt(task_1['urgency'])
		- parseInt(task_1['difficulty']);
	var task_2_total
		= parseInt(task_2['short_term'])
		+ parseInt(task_2['long_term'])
		+ parseInt(task_2['urgency'])
		- parseInt(task_2['difficulty']);
	return task_2_total - task_1_total;
}

function append_task(task, id, short_term, long_term, urgency, difficulty) {
	var task_template = $('.task.template')[0];
	var customize_template = function(clone) {
		$(clone).submit(ajax_form_submission);
		$(clone).addClass('task-id-' + id);
		$(clone).find('a.task')
			.html(Belt.escapeHTML(task))
			.attr('href', '/#' + id);
		$(clone).find('input.id').val(id);
		$(clone).find('.short-term').html(short_term);
		$(clone).find('.long-term').html(long_term);
		$(clone).find('.urgency').html(urgency);
		$(clone).find('.difficulty').html(difficulty);
		var total = parseInt(short_term) + parseInt(long_term) + parseInt(urgency) - parseInt(difficulty);
		$(clone).find('.total').html(total);
	};
	use_template(task_template, '.tasks', customize_template);
}



success_response_callbacks['new-task-success'] = new_task_success;

function new_task_success(post_data, new_task_data) {
	insert_new_task(
		post_data['task'],
		new_task_data,
		post_data['short-term'],
		post_data['long-term'],
		post_data['urgency'],
		post_data['difficulty']
	);
	display_tasks();
	$('.new-task input.task').val('');
	$('.new-task input.short-term, .new-task input.long-term, .new-task input.urgency, .new-task input.difficulty').val(5);
	$('.new-task input.task').select();
}

function insert_new_task(task, id, short_term, long_term, urgency, difficulty) {
	var new_task = {
		'task': decodeURIComponent(task),
		'id': id,
		'short_term': short_term,
		'long_term': long_term,
		'urgency': urgency,
		'difficulty': difficulty
	};
	tasks.push(new_task);
}



success_response_callbacks['delete-task-success'] = delete_task_success;

function delete_task_success(post_data, delete_task_data) {
	$('.task-id-' + post_data['id']).remove();
	for (var i in tasks) {
		if (tasks[i]['id'] === post_data['id']) {
			tasks.splice(i, 1);
			break;
		}
	}
}



function get_parent_id_from_url_hash() {
	var url_hash = window.location.hash.substring(1, window.location.hash.length);
	if (url_hash.length === 0) {return undefined;}
	var parent_id = parseInt(url_hash);
	if (parent_id.toString() !== url_hash) {return undefined;}
	return parent_id;
}
