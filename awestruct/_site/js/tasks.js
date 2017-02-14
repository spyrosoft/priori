function reset_tasks() {
	$('.tasks').html('');
	api_call(
		{'action': 'tasks'},
		'populate-tasks'
	);
}

success_response_callbacks['populate-tasks'] = populate_tasks_callback;

function populate_tasks_callback(post_data, tasks_data) {
	for (var i = 0; i < tasks_data.length; i++) {
		add_task(tasks_data[i]['task'], tasks_data[i]['id']);
	}
}

function add_task(task, id) {
	var task_template = $('.task.template')[0];
	$(task_template).find('a')
		.html(task)
		.attr('href', '/#' + id);
	use_template(task_template, '.tasks');
}



success_response_callbacks['new-task-success'] = new_task_success;

function new_task_success(post_data, new_task_data) {
	add_task(new_task_data['task'], new_task_data['id']);
	$('new-task').first().val('');
}



$('.new-task').select();
reset_tasks();