function reset_tasks() {
	$('.tasks').html('');
	api_call(
		{'action': 'tasks'},
		'populate-tasks'
	);
}

success_response_callbacks['populate-tasks'] = populate_tasks_callback;

function populate_tasks_callback(post_data, tasks_data) {
	console.log(tasks_data)
	for (var i in tasks_data) {
		add_task(
			tasks_data[i]['task'],
			tasks_data[i]['id'],
			tasks_data[i]['short-term'],
			tasks_data[i]['long-term'],
			tasks_data[i]['urgency'],
			tasks_data[i]['difficulty']
		);
	}
}

function add_task(task, id, short_term, long_term, urgency, difficulty) {
	var task_template = $('.task.template')[0];
	var customize_template = function(clone) {
		$(clone).submit(ajax_form_submission);
		$(clone).addClass('task-id-' + id);
		$(clone).find('a.task')
			.html(task)
			.attr('href', '/#' + id);
		$(clone).find('input.id').val(id);
		$(clone).find('.short-term')
			.html(short_term + ' + ');
		$(clone).find('.long-term')
			.html(long_term + ' + ');
		$(clone).find('.urgency')
			.html(urgency + ' - ');
		$(clone).find('.difficulty')
			.html(difficulty + ' = ');
		var total = parseInt(short_term) + parseInt(long_term) + parseInt(urgency) - parseInt(difficulty);
		$(clone).find('.total').html(total);
	};
	use_template(task_template, '.tasks', customize_template);
}



success_response_callbacks['new-task-success'] = new_task_success;

function new_task_success(post_data, new_task_data) {
	add_task(
		post_data['task'],
		new_task_data,
		post_data['short-term'],
		post_data['long-term'],
		post_data['urgency'],
		post_data['difficulty']
	);
	$('.new-task input.task, .new-task input.short-term, .new-task input.long-term, .new-task input.urgency, .new-task input.difficulty').val('');
	$('.new-task input.task').select();
}



success_response_callbacks['delete-task-success'] = delete_task_success;

function delete_task_success(post_data, delete_task_data) {
	$('.task-id-' + post_data['id']).remove();
}



$('.new-task input.task').select();
reset_tasks();