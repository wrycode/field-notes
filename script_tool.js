var lcode = "en_US";
var builtin_script_str = "teen_script.svg";
var custom_script_str;

var options = {
    "Debug": false,
    "Draw_bounding_boxes": false,
    "Form_stroke_width": 0.2,
    "Custom_script_svg_value":  "",    
    "Builtin_script_name": "teen_script.svg",
    "Language_code": "en_US",
    "Image_width": 500,
    "Input_text": "Let's start a\n new paragraph.",
    "Space_between_metaforms": 10,
    "Space_between_lines": 15,
    "Margin": 15
};

// Function to update options from form
function updateOptions() {
    options.Debug = document.getElementById('Debug').checked;
    options.Draw_bounding_boxes = document.getElementById('Draw_bounding_boxes').checked;
    options.Form_stroke_width = parseFloat(document.getElementById('Form_stroke_width').value);
    options.Builtin_script_name = document.getElementById('Builtin_script_name').value;
    options.Language_code = document.getElementById('Language_code').value;
    options.Image_width = parseInt(document.getElementById('Image_width').value, 10);
    options.Input_text = document.getElementById('Input_text').value;
    options.Space_between_metaforms = parseInt(document.getElementById('Space_between_metaforms').value, 10);
    options.Space_between_lines = parseInt(document.getElementById('Space_between_lines').value, 10);
    options.Margin = parseInt(document.getElementById('Margin').value, 10);
    // Potentially update display or storage with new options
    console.log(options); // Just for demonstration
}

// Set default values (if not using HTML inline values)
function setDefaultValues() {
    document.getElementById('Debug').checked = options.Debug;
    document.getElementById('Draw_bounding_boxes').checked = options.Draw_bounding_boxes;
    document.getElementById('Form_stroke_width').value = options.Form_stroke_width;
    document.getElementById('Builtin_script_name').value = options.Builtin_script_name;
    document.getElementById('Language_code').value = options.Language_code;
    document.getElementById('Image_width').value = options.Image_width;
    document.getElementById('Input_text').value = options.Input_text;
    document.getElementById('Space_between_metaforms').value = options.Space_between_metaforms;
    document.getElementById('Space_between_lines').value = options.Space_between_lines;
    document.getElementById('Margin').value = options.Margin;
}

// Initialize form with default values
setDefaultValues();

// Add event listener to form elements to update options object
document.getElementById("form-container").addEventListener('change', updateOptions);

document.getElementById('renderButton').addEventListener('click', function() {
    var debug_output = document.getElementById('render_text_output');
    
    if (builtin_script_str) {
	Render(JSON.stringify(options))
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
var languages = [
    { code: "ar", name: "Arabic" },
    { code: "de", name: "German" },
    { code: "en_UK", name: "English (UK)" },
    { code: "en_US", name: "English (US)" },
    { code: "eo", name: "Esperanto" },
    { code: "es_ES", name: "Spanish (Spain)" },
    { code: "es_MX", name: "Spanish (Mexico)" },
    { code: "fa", name: "Persian" },
    { code: "fi", name: "Finnish" },
    { code: "fr_FR", name: "French (France)" },
    { code: "fr_QC", name: "French (Quebec)" },
    { code: "ja", name: "Japanese" },
    { code: "jam", name: "Jamaican Patois" },
    { code: "ma", name: "Moroccan Arabic" },
    { code: "nb", name: "Norwegian Bokmål" },
    { code: "or", name: "Oriya" },
    { code: "sv", name: "Swedish" },
    { code: "sw", name: "Swahili" },
    { code: "vi_C", name: "Vietnamese (Central)" },
    { code: "vi_N", name: "Vietnamese (Northern)" },
    { code: "vi_S", name: "Vietnamese (Southern)" },
    { code: "yue", name: "Cantonese" },
    { code: "zh_hans", name: "Chinese (Simplified)" },
    { code: "zh_hant", name: "Chinese (Traditional)" }
];

const logTextarea = document.getElementById('log');
logTextarea.addEventListener('input', function() {
    this.style.height = 'auto';
    this.style.height = this.scrollHeight + 'px';
});

document.addEventListener('DOMContentLoaded', function () {
  var languageSelect = document.getElementById('Language_code');

  // Function to populate language options
  function populateLanguages() {
    languages.forEach(function(language) {
      var option = document.createElement('option');
      option.value = language.code;
      option.text = language.name;
      languageSelect.appendChild(option);
    });
  }

  // Populate languages on document load
  populateLanguages();
  
  // Rest of your existing code to handle form updates, etc.
});



// let HF_intro_text = `YOU don't know about me without you have read a book by the name of The Adventures of Tom Sawyer; but that ain't no matter.  That book was made by Mr. Mark Twain, and he told the truth, mainly.  There was things which he stretched, but mainly he told the truth.  That is nothing.  I never seen anybody but lied one time or another, without it was Aunt Polly, or the widow, or maybe Mary.  Aunt Polly—Tom's Aunt Polly, she is—and Mary, and the Widow Douglas is all told about in that book, which is mostly a true book, with some stretchers, as I said before.`;

// document.getElementById("Input_text").value = HF_intro_text;
