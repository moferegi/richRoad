export const timeFormat = (time, type = 'YYYY-MM-DD', line = "-") => {
	if (!time) return ''
	let date = new Date(time);
	let Y = date.getFullYear() + line;
	let M = (date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1) + line;
	let D = (date.getDate() < 10 ? '0' + date.getDate() : date.getDate());
	let h = (date.getHours() < 10 ? '0' + date.getHours() : date.getHours()) + ':';
	let m = (date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()) + ':';
	let s = (date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds());
	if (type == 'YYYY-MM-DD hh-mm-ss') {
		return Y + M + D + ' ' + h + m + s
	}
	if (type == 'YYYY-MM-DD') {
		return Y + M + D
	} else {
		return h + m + s
	}
}