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
		add_task(
			tasks_data[i]['task'],
			tasks_data[i]['id'],
			tasks_data[i]['difficulty'],
			tasks_data[i]['short_term'],
			tasks_data[i]['long_term']
		);
	}
}

function add_task(task, id, difficulty, short_term, long_term) {
	var task_template = $('.task.template')[0];
	var customize_template = function(clone) {
		$(clone).submit(ajax_form_submission);
		$(clone).find('a.task')
			.html(task)
			.attr('href', '/#' + id);
		$(clone).find('input.id').val(id);
		$(clone).find('.difficulty')
			.html(difficulty);
		$(clone).find('.short-term')
			.html(short_term);
		$(clone).find('.long-term')
			.html(long_term);
	};
	use_template(task_template, '.tasks', customize_template);
}



success_response_callbacks['new-task-success'] = new_task_success;

function new_task_success(post_data, new_task_data) {
	add_task(
		post_data['task'],
		new_task_data,
		post_data['difficulty'],
		post_data['short-term'],
		post_data['long-term']
	);
	$('new-task').first().val('');
}



$('.new-task').select();
reset_tasks();