// DO NOT EDIT. This file is automatically generated.

export interface ResponsePagination {
	offset: number;
	size: number;
	total: number;
}

export interface RedoTaskRequest {
	id: number;
}

export interface DeleteTaskRequest {
	id: number;
}

export interface CommonCols {
	ts_added: number;
	ts_started: number;
	ts_completed: number;
	username: string;
	failure: boolean;
	repeat_after_days: number;
	swarming_logs: string;
	task_done: boolean;
	swarming_task_id: string;
	id: number;
	can_redo: boolean;
	can_delete: boolean;
	future_date: boolean;
	task_type: string;
	get_url: string;
	delete_url: string;
}

export interface Permissions {
	DeleteAllowed: boolean;
	RedoAllowed: boolean;
}

export interface GetTasksResponse {
	data: any;
	permissions: Permissions[] | null;
	pagination: ResponsePagination | null;
	ids: number[] | null;
}

export interface BenchmarksPlatformsResponse {
	benchmarks: { [key: string]: string };
	platforms: { [key: string]: string };
}

export interface TaskPrioritiesResponse {
	task_priorities: { [key: number]: string };
}

export interface PageSet {
	key: string;
	description: string;
}

export interface CLDataResponse {
	cl: string;
	subject: string;
	url: string;
	modified: string;
	chromium_patch: string;
	skia_patch: string;
	v8_patch: string;
	catapult_patch: string;
}

export interface CompletedTask {
	type: string;
	username: string;
	description: string;
	ts_completed: number;
}

export interface CompletedTaskResponse {
	unique_users: number;
	completed_tasks: CompletedTask[] | null;
}

export interface AdminDatastoreTask {
	ts_added: number;
	ts_started: number;
	ts_completed: number;
	username: string;
	failure: boolean;
	repeat_after_days: number;
	swarming_logs: string;
	task_done: boolean;
	swarming_task_id: string;
	id: number;
	can_redo: boolean;
	can_delete: boolean;
	future_date: boolean;
	task_type: string;
	get_url: string;
	delete_url: string;
	page_sets: string;
	is_test_page_set: boolean;
}

export interface AdminAddTaskVars {
	username: string;
	ts_added: string;
	repeat_after_days: string;
	page_sets: string;
}

export interface ChromiumAnalysisDatastoreTask {
	ts_added: number;
	ts_started: number;
	ts_completed: number;
	username: string;
	failure: boolean;
	repeat_after_days: number;
	swarming_logs: string;
	task_done: boolean;
	swarming_task_id: string;
	id: number;
	can_redo: boolean;
	can_delete: boolean;
	future_date: boolean;
	task_type: string;
	get_url: string;
	delete_url: string;
	benchmark: string;
	page_sets: string;
	is_test_page_set: boolean;
	benchmark_args: string;
	browser_args: string;
	description: string;
	custom_webpages_gspath: string;
	chromium_patch_gspath: string;
	skia_patch_gspath: string;
	catapult_patch_gspath: string;
	benchmark_patch_gspath: string;
	v8_patch_gspath: string;
	run_in_parallel: boolean;
	platform: string;
	run_on_gce: boolean;
	raw_output: string;
	value_column_name: string;
	match_stdout_txt: string;
	chromium_hash: string;
	apk_gspath: string;
	telemetry_isolate_hash: string;
	cc_list: string[] | null;
	task_priority: number;
	group_name: string;
}

export interface ChromiumAnalysisAddTaskVars {
	username: string;
	ts_added: string;
	repeat_after_days: string;
	benchmark: string;
	page_sets: string;
	custom_webpages: string;
	benchmark_args: string;
	browser_args: string;
	desc: string;
	chromium_patch: string;
	skia_patch: string;
	catapult_patch: string;
	benchmark_patch: string;
	v8_patch: string;
	run_in_parallel: boolean;
	platform: string;
	run_on_gce: boolean;
	value_column_name: string;
	match_stdout_txt: string;
	chromium_hash: string;
	apk_gs_path: string;
	telemetry_isolate_hash: string;
	cc_list: string[] | null;
	task_priority: string;
	group_name: string;
}

export interface ChromiumPerfDatastoreTask {
	ts_added: number;
	ts_started: number;
	ts_completed: number;
	username: string;
	failure: boolean;
	repeat_after_days: number;
	swarming_logs: string;
	task_done: boolean;
	swarming_task_id: string;
	id: number;
	can_redo: boolean;
	can_delete: boolean;
	future_date: boolean;
	task_type: string;
	get_url: string;
	delete_url: string;
	benchmark: string;
	platform: string;
	run_on_gce: boolean;
	page_sets: string;
	is_test_page_set: boolean;
	repeat_runs: number;
	run_in_parallel: boolean;
	benchmark_args: string;
	browser_args_no_patch: string;
	browser_args_with_patch: string;
	description: string;
	custom_webpages_gspath: string;
	chromium_patch_gspath: string;
	blink_patch_gspath: string;
	skia_patch_gspath: string;
	catapult_patch_gspath: string;
	benchmark_patch_gspath: string;
	chromium_patch_base_build_gspath: string;
	v8_patch_gspath: string;
	results: string;
	no_patch_raw_output: string;
	with_patch_raw_output: string;
	chromium_hash: string;
	cc_list: string[] | null;
	task_priority: number;
	group_name: string;
	value_column_name: string;
}

export interface ChromiumPerfAddTaskVars {
	username: string;
	ts_added: string;
	repeat_after_days: string;
	benchmark: string;
	platform: string;
	run_on_gce: string;
	page_sets: string;
	custom_webpages: string;
	repeat_runs: string;
	run_in_parallel: string;
	benchmark_args: string;
	browser_args_nopatch: string;
	browser_args_withpatch: string;
	desc: string;
	chromium_hash: string;
	cc_list: string[] | null;
	task_priority: string;
	group_name: string;
	value_column_name: string;
	chromium_patch: string;
	blink_patch: string;
	skia_patch: string;
	catapult_patch: string;
	benchmark_patch: string;
	v8_patch: string;
	chromium_patch_base_build: string;
}

export interface MetricsAnalysisDatastoreTask {
	ts_added: number;
	ts_started: number;
	ts_completed: number;
	username: string;
	failure: boolean;
	repeat_after_days: number;
	swarming_logs: string;
	task_done: boolean;
	swarming_task_id: string;
	id: number;
	can_redo: boolean;
	can_delete: boolean;
	future_date: boolean;
	task_type: string;
	get_url: string;
	delete_url: string;
	metric_name: string;
	analysis_task_id: string;
	analysis_output_link: string;
	benchmark_args: string;
	description: string;
	custom_traces_gspath: string;
	chromium_patch_gspath: string;
	catapult_patch_gspath: string;
	raw_output: string;
	value_column_name: string;
	cc_list: string[] | null;
	task_priority: number;
}

export interface MetricsAnalysisAddTaskVars {
	username: string;
	ts_added: string;
	repeat_after_days: string;
	metric_name: string;
	custom_traces: string;
	analysis_task_id: string;
	analysis_output_link: string;
	benchmark_args: string;
	desc: string;
	chromium_patch: string;
	catapult_patch: string;
	value_column_name: string;
	cc_list: string[] | null;
	task_priority: string;
}