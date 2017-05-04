/* -------------------- Globals -------------------- */
var all_tasks = {};
var task_template = {
	'name' : ''
	, 'weights-values' : []
	, 'children' : []
	, 'children-weights' : []
	, 'parent-id' : ''
};
var weights_template = {
	'Short Term Gain': true
	, 'Long Term Gain': true
	, 'Difficulty': false
};
/* -------------------- Globals -------------------- */


/* -------------------- Button Click Events -------------------- */
$('button.edit-weights').on('click', edit_weights);
function edit_weights() {
	if ($('#edit-weights').is(':visible')) {
		$('#edit-weights').hide();
	} else {
		$('#edit-weights').show();
	}
}

$('#edit-weights button.add').on('click', add_weight);
function add_weight() {
	
}
/* -------------------- Button Click Events -------------------- */


/* -------------------- Page Load -------------------- */
$('.new-task input.task').select();
populate_tasks();
/* -------------------- Page Load -------------------- */


function populate_tasks() {
	var post_data = {'action': 'get-tasks'};
	api_call(post_data);
}

success_response_callbacks['get-tasks'] = function(post_data, tasks_data) {
	if (tasks_data.length === 0) {
		all_tasks['-'] = default_task();
	} else if (tasks_data['-'] === undefined) {
		display_error('An error occurred while retrieving your tasks. An admin has been notified.');
		notify_admin('The tasks received from the server didn\' have a top task: ' + JSON.stringify(tasks_data));
		return;
	}
	display_tasks(all_tasks['-']);
};

function default_task() {
	var task = clone_object(task_template);
	task['name'] = 'Achievements';
	task['children-weights'] = default_weights();
	return task;
}

function default_weights() {
	return clone_object(weights_template);
}

function display_tasks(task) {
	$('.tasks').html('');
	sort_children_tasks(task);
	for (var c in task['children']) {
		append_task(task['children'][c]);
	}
	$('.weights').html('');
	append_weights(task['children-weights']);
}

function sort_children_tasks(task) {
	if (task['children'].length === 0) { return; }
	invert_bad_weights_signs(task);
	task['children'].sort(sort_tasks);
	invert_bad_weights_signs(task);
}

function invert_bad_weights_signs(task) {
	for (var w in task['children-weights']) {
		if (!task['children-weights'][w]) {
			for (var c in task['children']) {
				task['children'][c]['weights'][w] *= -1;
			}
		}
	}
}

function sort_tasks(task_1, task_2) {
	var task_1_total = 0;
	for (var i1 in task_1['weights']) {
		task_1_total += task_1['weights'][i1];
	}
	var task_2_total = 0;
	for (var i2 in task_1['weights']) {
		task_1_total += task_1['weights'][i2];
	}
	return task_2_total - task_1_total;
}

function append_task(task) {
	var task_template = $('.task.template')[0];
	var closure = function(clone) { customize_task_template(clone, task); };
	use_template(task_template, '.tasks', closure);
}

function customize_task_template(task_clone_element, task) {
	var task_clone = $(task_clone_element);
	task_clone.submit(update_task);
	task_clone.attr('id', 'task-id-' + task['id']);
	task_clone.find('a.task')
		.html(Belt.escapeHTML(task))
		.on('click', child_task_clicked);
	for (var w in task['children-weights']) {
		var task_weight_template = $('.task-weight.template')[0];
		var closure = function(clone) { customize_task_weights_template(clone, task['children-weights'], task['children-weights'][w]); };
		use_template(task_template, '.tasks', closure);
	}
	//TODO: Delete Me:
	//task_clone.find('.difficulty').html(difficulty);
	//var total = parseInt(short_term) + parseInt(long_term) + parseInt(urgency) - parseInt(difficulty);
	//task_clone.find('.total').html(total);
}

function customize_task_weights_template(weight_clone_element, weights, weight) {
}

function append_weights(weights) {
	var weight_template = $('.weight.template')[0];
	for (var w in weights) {
		var closure = function(clone) { customize_weight_template(clone, weights, weights[w]); };
	}
	use_template(weight_template, '.weights', closure);
}

function customize_weight_template(weight_clone_element, weights, weight) {
	var weight_clone = $(weight_clone_element);
	if (weights[weight]) {
		weight_clone.find('.good-or-bad')
			.addClass('good')
			.html('+');
	} else {
		weight_clone.find('.good-or-bad')
			.addClass('bad')
			.html('-');
	}
}

success_response_callbacks['new-task'] = new_task_success;

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

success_response_callbacks['delete-task'] = delete_task_success;

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

function clone_object(object_to_clone) {
	return $.extend(true, {}, object_to_clone);
}
