// DO NOT EDIT. This file is automatically generated.

export interface Alert {
	id: number;
	display_name: string;
	query: string;
	alert: string;
	interesting: number;
	bug_uri_template: string;
	algo: RegressionDetectionGrouping;
	step: StepDetection;
	state: ConfigState;
	owner: string;
	step_up_only: boolean;
	direction: Direction;
	radius: number;
	k: number;
	group_by: string;
	sparse: boolean;
	minimum_num: number;
	category: string;
}

export interface CommitDetail {
	offset: CommitNumber;
	author: string;
	message: string;
	url: string;
	hash: string;
	ts: number;
}

export interface ValuePercent {
	value: string;
	percent: number;
}

export interface StepFit {
	least_squares: number;
	turning_point: number;
	step_size: number;
	regression: number;
	status: StepFitStatus;
}

export interface ColumnHeader {
	offset: CommitNumber;
	timestamp: number;
}

export interface ClusterSummary {
	centroid: number[] | null;
	shortcut: string;
	param_summaries2: ValuePercent[] | null;
	step_fit: StepFit | null;
	step_point: ColumnHeader | null;
	num: number;
}

export interface FrameRequest {
	begin: number;
	end: number;
	formulas: string[] | null;
	queries: string[] | null;
	hidden: string[] | null;
	keys: string;
	tz: string;
	num_commits: number;
	request_type: RequestType;
}

export interface DataFrame {
	traceset: TraceSet;
	header: (ColumnHeader | null)[] | null;
	paramset: ParamSet;
	skip: number;
}

export interface FrameResponse {
	dataframe: DataFrame | null;
	skps: number[] | null;
	msg: string;
}

export interface TriageStatus {
	status: Status;
	message: string;
}

export interface Regression {
	low: ClusterSummary | null;
	high: ClusterSummary | null;
	frame: FrameResponse | null;
	low_status: TriageStatus;
	high_status: TriageStatus;
}

export interface RegressionRow {
	cid: CommitDetail | null;
	regression: Regression | null;
}

export interface DryRunStatus {
	finished: boolean;
	message: string;
	regressions: (RegressionRow | null)[] | null;
}

export interface UIDomain {
	begin: number;
	end: number;
	num_commits: number;
	request_type: RequestType;
}

export interface StartDryRunRequest {
	config: Alert;
	domain: UIDomain;
}

export interface StartDryRunResponse {
	id: string;
}

export interface ClusterStartResponse {
	id: string;
}

export interface ClusterSummaries {
	Clusters: (ClusterSummary | null)[] | null;
	StdDevThreshold: number;
	K: number;
}

export interface RegressionDetectionResponse {
	summary: ClusterSummaries | null;
	frame: FrameResponse | null;
}

export interface ClusterStatus {
	state: ProcessState;
	message: string;
	value: RegressionDetectionResponse | null;
}

export interface CommitID {
	offset: CommitNumber;
}

export interface CommitDetailsRequest {
	cid: CommitID;
	traceid: string;
}

export interface CountHandlerRequest {
	q: string;
	begin: number;
	end: number;
}

export interface CountHandlerResponse {
	count: number;
	paramset: ParamSet;
}

export interface RangeRequest {
	offset: CommitNumber;
	begin: number;
	end: number;
}

export interface RegressionRangeRequest {
	begin: number;
	end: number;
	subset: Subset;
	alert_filter: string;
}

export interface RegressionRow {
	cid: CommitDetail | null;
	columns: (Regression | null)[] | null;
}

export interface RegressionRangeResponse {
	header: (Alert | null)[] | null;
	table: (RegressionRow | null)[] | null;
	categories: string[] | null;
}

export interface SkPerfConfig {
	radius: number;
	key_order: string[] | null;
	num_shift: number;
	interesting: number;
	step_up_only: boolean;
	commit_range_url: string;
	demo: boolean;
}

export interface TriageRequest {
	cid: CommitID | null;
	alert: Alert;
	triage: TriageStatus;
	cluster_type: string;
}

export interface TriageResponse {
	bug: string;
}

export interface TryBugRequest {
	bug_uri_template: string;
}

export interface TryBugResponse {
	url: string;
}

export interface Current {
	commit: CommitDetail | null;
	alert: Alert | null;
	message: string;
}

export interface FullSummary {
	summary: ClusterSummary;
	triage: TriageStatus;
	frame: FrameResponse;
}

export interface Domain {
	n: number;
	end: string;
	offset: number;
}

export interface RegressionDetectionRequest {
	alert: Alert | null;
	domain: Domain;
	query: string;
	step: number;
	total_queries: number;
}

export type RegressionDetectionGrouping = string;

export type StepDetection = "" | "absolute" | "percent" | "cohen";

export type ConfigState = "ACTIVE" | "DELETED";

export type Direction = "UP" | "DOWN" | "BOTH";

export type CommitNumber = number;

export type StepFitStatus = "Low" | "High" | "Uninteresting";

export type RequestType = 1 | 0;

export type Trace = number[] | null;

export type TraceSet = { [key: string]: Trace };

export type ParamSet = { [key: string]: string[] | null };

export type Status = "" | "positive" | "negative" | "untriaged";

export type ProcessState = "Running" | "Success" | "Error";

export type Subset = "all" | "regressions" | "untriaged";

export type ClusterAlgo = "kmeans" | "stepfit";
