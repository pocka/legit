// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

class TabWidthSelector {
	static LOCAL_STORAGE_KEY = "legit_tab_width";

	static #isValidWidth(width) {
		if (typeof width !== "string") {
			return false;
		}

		const n = parseInt(width, 10);

		return Number.isFinite(n) && n > 0;
	}

	static #selector() {
		return document.getElementById("global-tab-width");
	}

	static sync() {
		const saved = localStorage.getItem(this.LOCAL_STORAGE_KEY);
		const value = this.#isValidWidth(saved) ? saved : "4";
		document.body.style.tabSize = value;

		const selector = this.#selector();
		if (selector) {
			selector.value = saved;
		}
	}

	static addListener() {
		this.#selector().addEventListener("change", event => {
			if (!this.#isValidWidth(event.currentTarget.value)) {
				return;
			}

			localStorage.setItem(this.LOCAL_STORAGE_KEY, event.currentTarget.value);
			this.sync();
		});
	}
}

class ThemeColor {
	static LOCAL_STORAGE_KEY = "legit_theme_color";
	static VALID_COLORS = ["red", "orange", "green", "gray", "pink", "purple", "yellow"];

	static #selector() {
		return document.getElementById("theme-color");
	}

	static sync() {
		const value = localStorage.getItem(this.LOCAL_STORAGE_KEY) || "";
		if (value) {
			document.documentElement.dataset.themeColor = value;
		} else {
			document.documentElement.dataset.themeColor = "";
		}

		const selector = this.#selector();
		if (selector) {
			selector.value = value;
		}
	}

	static addListener() {
		this.#selector()?.addEventListener("change", event => {
			const value = event.currentTarget.value;
			if (value && !this.VALID_COLORS.includes(value)) {
				return;
			}

			if (value) {
				localStorage.setItem(this.LOCAL_STORAGE_KEY, value);
			} else {
				localStorage.removeItem(this.LOCAL_STORAGE_KEY);
			}
			this.sync();
		});
	}
}

class TreeLayout {
	static LOCAL_STORAGE_KEY = "legit_tree_layout";
	static VALID_LAYOUTS = ["compact"];

	static #selector() {
		return document.getElementById("tree-layout");
	}

	static sync() {
		const value = localStorage.getItem(this.LOCAL_STORAGE_KEY) || "";

		for (const tree of document.getElementsByClassName("tree-row")) {
			tree.dataset.layout = value;
		}

		const selector = this.#selector();
		if (selector) {
			selector.value = value;
		}
	}

	static addListener() {
		this.#selector()?.addEventListener("change", event => {
			const value = event.currentTarget.value;
			if (value && !this.VALID_LAYOUTS.includes(value)) {
				return;
			}

			if (value) {
				localStorage.setItem(this.LOCAL_STORAGE_KEY, value);
			} else {
				localStorage.removeItem(this.LOCAL_STORAGE_KEY);
			}
			this.sync();
		});
	}
}

class FileModeFormat {
	static LOCAL_STORAGE_KEY = "legit_file_mode_format";
	static VALID_FORMATS = ["octal", "hidden"];
	static FILE_MODE_PATTERN = /^[d-]([r-][w-][x-]){3}$/;

	static #selector() {
		return document.getElementById("file-mode-format");
	}

	static sync() {
		const value = localStorage.getItem(this.LOCAL_STORAGE_KEY) || "";

		for (const el of document.querySelectorAll("[data-file-mode]")) {
			const mode = el.dataset.fileMode;
			if (!this.FILE_MODE_PATTERN.test(mode || "")) {
				continue;
			}

			el.dataset.format = value;
			switch (value) {
				case "octal": {
					const n = 0
						+ (mode[1] === "r" ? 0o400 : 0)
						+ (mode[2] === "w" ? 0o200 : 0)
						+ (mode[3] === "x" ? 0o100 : 0)
						+ (mode[4] === "r" ? 0o40 : 0)
						+ (mode[5] === "w" ? 0o20 : 0)
						+ (mode[6] === "x" ? 0o10 : 0)
						+ (mode[7] === "r" ? 0o4 : 0)
						+ (mode[8] === "w" ? 0o2 : 0)
						+ (mode[9] === "x" ? 0o1 : 0);
					el.textContent = n.toString(8).padStart(3, "0");
					break;
				}
				default:
					el.textContent = mode;
					break;
			}
		}

		const selector = this.#selector();
		if (selector) {
			selector.value = value;
		}
	}

	static addListener() {
		this.#selector()?.addEventListener("change", event => {
			const value = event.currentTarget.value;
			if (value && !this.VALID_FORMATS.includes(value)) {
				return;
			}

			if (value) {
				localStorage.setItem(this.LOCAL_STORAGE_KEY, value);
			} else {
				localStorage.removeItem(this.LOCAL_STORAGE_KEY);
			}
			this.sync();
		});
	}
}

class FileSizeFormat {
	static LOCAL_STORAGE_KEY = "legit_file_size_format";
	static VALID_FORMATS = ["hidden"];

	static #selector() {
		return document.getElementById("file-size-format");
	}

	static sync() {
		const value = localStorage.getItem(this.LOCAL_STORAGE_KEY) || "";

		for (const el of document.querySelectorAll("[data-file-size]")) {
			const size = parseInt(el.dataset.fileSize, 10);
			if (!Number.isFinite(size)) {
				continue;
			}

			el.dataset.format = value;
		}

		const selector = this.#selector();
		if (selector) {
			selector.value = value;
		}
	}

	static addListener() {
		this.#selector()?.addEventListener("change", event => {
			const value = event.currentTarget.value;
			if (value && !this.VALID_FORMATS.includes(value)) {
				return;
			}

			if (value) {
				localStorage.setItem(this.LOCAL_STORAGE_KEY, value);
			} else {
				localStorage.removeItem(this.LOCAL_STORAGE_KEY);
			}
			this.sync();
		});
	}
}

function setupPreferencesDialog() {
	const preferences = document.getElementById("preferences-dialog");
	document.getElementById("preferences-dialog-trigger")?.addEventListener("click", (event) => {
		preferences?.showModal();
	});
	document.getElementById("preferences-dialog-closer")?.addEventListener("click", (event) => {
		preferences?.close();
	});
}

TabWidthSelector.addListener();
ThemeColor.addListener();
TreeLayout.addListener();
FileModeFormat.addListener();
FileSizeFormat.addListener();

setupPreferencesDialog();

window.addEventListener("pageshow", event => {
	TabWidthSelector.sync();
	ThemeColor.sync();
	TreeLayout.sync();
	FileModeFormat.sync();
	FileSizeFormat.sync();
});
