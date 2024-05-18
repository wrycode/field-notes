var lcode = "en_US";
var builtin_script_str = "teen_script.svg";
var custom_script_str;

document.getElementById("custom_script_select").addEventListener("change", function(e){
    if(e.target.files.length > 0) {
        builtin_script_str = null;
        custom_script = e.target.files[0];
        document.getElementById("builtin_script_select").disabled = true;
    } else {
        document.getElementById("builtin_script_select").disabled = false;
    }
});


document.getElementById('renderButton').addEventListener('click', function() {
    var debug_output = document.getElementById('render_text_output');
    
    if (builtin_script_str) {
	renderSVG(false, builtin_script_str, lcode)
    } else {
	var file = document.getElementById('custom_script_select').files[0];
	var reader = new FileReader();
	reader.readAsText(file, 'UTF-8');
        reader.onload = function (evt) {
            renderSVG(true, evt.target.result, lcode);
        }
    reader.onerror = function (evt) {
        debug_output.textContent = "An error occurred reading the file.";
        }
    }
});

var languages =
    ["ar", "de", "en_UK", "en_US", "eo", "es_ES", "es_MX", "fa",
     "fi", "fr_FR", "fr_QC", "ja", "jam", "ma", "nb", "or", "sv",
     "sw","vi_C", "vi_N", "vi_S", "yue", "zh_hans", "zh_hant"]

var lang_select = document.getElementById('lang_select');

for(var i = 0; i < languages.length; i++) {
    var opt = languages[i];
    var el = document.createElement("option");
    el.textContent = opt;
    // el.value = opt;
    lang_select.appendChild(el);
}

lang_select.onchange = function() {
    lcode = lang_select.value;
};

lang_select.value = "en_US"
