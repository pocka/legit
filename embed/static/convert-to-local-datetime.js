// Copyright 2025 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT
//
// ===
//
// This script scans every "<time>" elements with "data-local-format" attribute,
// replaces each element's text node with local datetime.
// If a "<time>" element does not have valid "datetime" attribute, this script
// skips the "<time>" element.
//
// "data-local-format" attribute accepts one of the below:
// - "datetime" ... displays date and time (default)
// - "date"     ... displays date
// - "time"     ... displays time
// - "relative" ... displays relative datetime (e.g. 3 days ago)
//
// This script intentionally uses scratch serialization instead of locale-aware
// APIs; locale strings often includes non-English characters thus does not
// incorporate with the rest of a page. As majority of readers are familiar with
// softwares, using RFC3339-ish style date and time format fits better than an
// uncontrollable locale string.

const MS_PER_MIN = 60 * 1000;
const MS_PER_HOUR = MS_PER_MIN * 60;
const MS_PER_DAY = MS_PER_HOUR * 24;
const MS_PER_WEEK = MS_PER_DAY * 7;
const MS_PER_MONTH = MS_PER_DAY * 30;
const MS_PER_YEAR = MS_PER_MONTH * 12;

function toTimeString(date) {
	return date.getHours().toString().padStart(2, "0") + ":"
		+ date.getMinutes().toString().padStart(2, "0") + ":"
		+ date.getSeconds().toString().padStart(2, "0");
}

function toDateString(date) {
	// I don't expect this software/code will survive for 8000yrs.
	return date.getFullYear().toString().padStart(4, "0") + "-"
		// stupid spec made by American
		+ (date.getMonth() + 1).toString().padStart(2, "0") + "-"
		+ date.getDate().toString().padStart(2, "0");
}

function toDateTimeString(date) {
	return toDateString(date) + " " + toTimeString(date);
}

function toRelativeDatetimeString(date, now = Date.now()) {
	if (typeof Intl.RelativeTimeFormat !== "function") {
		return toDateTimeString(date);
	}

	const formatter = new Intl.RelativeTimeFormat(document.documentElement.lang);
	const delta = now - date.valueOf();

	if (delta < MS_PER_MIN) {
		return formatter.format(-Math.floor(delta / 1000), "second");
	}

	if (delta < MS_PER_HOUR) {
		return formatter.format(-Math.floor(delta / MS_PER_MIN), "minute");
	}

	if (delta < MS_PER_DAY) {
		return formatter.format(-Math.floor(delta / MS_PER_HOUR), "hour");
	}

	if (delta < MS_PER_WEEK) {
		return formatter.format(-Math.floor(delta / MS_PER_DAY), "day");
	}

	if (delta < MS_PER_MONTH) {
		return formatter.format(-Math.floor(delta / MS_PER_WEEK), "week");
	}

	if (delta < MS_PER_YEAR) {
		return formatter.format(-Math.floor(delta / MS_PER_MONTH), "month");
	}

	return formatter.format(-Math.floor(delta / MS_PER_YEAR), "year");
}

const targets = document.querySelectorAll("time[data-local-format]");
for (const target of targets) {
	const datetime = target.dateTime && Date.parse(target.dateTime);
	if (Number.isNaN(datetime)) {
		continue;
	}

	const originalText = target.textContent.trim();

	switch (target.dataset.localFormat || "datetime") {
		case "datetime":
			target.textContent = toDateTimeString(new Date(datetime));
			break;
		case "date":
			target.textContent = toDateString(new Date(datetime));
			break;
		case "time":
			target.textContent = toTimeString(new Date(datetime));
			break;
		case "relative":
			target.textContent = toRelativeDatetimeString(new Date(datetime));
			break;
		default:
			console.warn(
				`Found unknown value on data-local-format attribute: "${target.dataset.localFormat}".\n`
					+ `Available values are: "datetime", "date", "time"\n`
					+ "Skipping text replacement.",
			);
			break;
	}

	target.title = originalText;
}
