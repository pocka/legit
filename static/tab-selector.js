// Copyright 2025 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

const LOCAL_STORAGE_KEY = "legit_tab_width";

function isValidTabWidth(width) {
	if (typeof width !== "string") {
		return false;
	}

	const n = parseInt(width, 10);

	return Number.isFinite(n) && n > 0;
}

function applyTabWidth(width) {
	localStorage.setItem(LOCAL_STORAGE_KEY, width);
	document.body.style.tabSize = width;
}

const selector = document.getElementById("global_tab_width");
if (selector) {
	const saved = localStorage.getItem(LOCAL_STORAGE_KEY);
	const value = isValidTabWidth(saved) ? saved : "4";

	selector.value = value;
	applyTabWidth(value);

	selector.addEventListener("change", event => {
		if (!isValidTabWidth(event.currentTarget.value)) {
			return;
		}

		applyTabWidth(event.currentTarget.value);
	});

	selector.disabled = false;
}
