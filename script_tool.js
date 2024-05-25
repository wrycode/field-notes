var defaults = {
    "Debug": false,
    "Draw_bounding_boxes": false,
    "Form_stroke_width": 0.35,
    "Custom_script_svg_value":  "",    
    "Builtin_script_name": "teen_script.svg",
    "Language_code": "en_US",
    "Input_text": "Let's start a\nnew paragraph. Here's a bunch of random text. How does it look?",
    "Space_between_metaforms": 5,
    "Space_between_lines": 10,
    "Margin": 5,
}

// copy defaults to options variable
var options = Object.assign({}, defaults);

function updateOptions() {
    options.Debug = document.getElementById('Debug').checked;
    options.Draw_bounding_boxes = document.getElementById('Draw_bounding_boxes').checked;
    options.Form_stroke_width = parseFloat(document.getElementById('Form_stroke_width').value);
    options.Builtin_script_name = document.getElementById('Builtin_script_name').value;
    options.Language_code = document.getElementById('Language_code').value;
    options.Input_text = document.getElementById('Input_text').value;
    options.Space_between_metaforms = parseInt(document.getElementById('Space_between_metaforms').value, 10);
    options.Space_between_lines = parseInt(document.getElementById('Space_between_lines').value, 10);
    options.Margin = parseInt(document.getElementById('Margin').value, 10);
    console.log(options);
}

function setDefaultValues() {
    document.getElementById('Debug').checked = defaults.Debug;
    document.getElementById('Draw_bounding_boxes').checked = defaults.Draw_bounding_boxes;
    document.getElementById('Form_stroke_width').value = defaults.Form_stroke_width;
    document.getElementById('Builtin_script_name').value = defaults.Builtin_script_name;
    document.getElementById('Language_code').value = defaults.Language_code;

    range_image_width_input.value = 800;
    // Create a new 'change' event
    var event = new Event('input');
    // Dispatch it on the input element
    range_image_width_input.dispatchEvent(event);
    drawBackground("yellow");

    document.getElementById('Input_text').value = defaults.Input_text;
    document.getElementById('Space_between_metaforms').value = defaults.Space_between_metaforms;
    document.getElementById('Space_between_lines').value = defaults.Space_between_lines;
    document.getElementById('Margin').value = defaults.Margin;
    options = Object.assign({}, defaults);    
}

// document.getElementById('Image_width').value = 500
document.getElementById('defaultsButton').addEventListener('click', function() {
    setDefaultValues()
});

    // Access the input element
var range_image_width_input = document.getElementById("range_image_width_input");
    // Access the canvas element
var canvas = document.getElementById("canvas");
var ctx = canvas.getContext("2d");


function drawBackground(color) {
      ctx.fillStyle = color; // Set color to fill
      ctx.fillRect(0, 0, canvas.width, canvas.height); // Draws the filled rectangle
}

function updateSlider(value) {
    document.getElementById('range_image_width_input').value = value;
}

function updateNumber(value) {
    document.getElementById('number_image_width_input').value = value;
}

function updateBoth(value) {
    updateSlider(value);
    updateNumber(value);
}

const image_with_input_elements = [
  document.querySelector('#range_image_width_input'),
  document.querySelector('#number_image_width_input')
];

image_with_input_elements.forEach(function(elem) {
    elem.addEventListener("input", function() {
    var newWidth = this.value; 
    var newHeight = this.value;

    canvas.width = newWidth;
    canvas.height = newHeight;

      drawBackground("yellow");
    });
});

number_image_width_input.value = 800;
// Create a new 'change' event
var event = new Event('input');
// Dispatch it on the input element
range_image_width_input.dispatchEvent(event);
drawBackground("yellow");


document.getElementById('Custom_script_svg_value').addEventListener('change', function(event) {
  var file = event.target.files[0];
  if (file) {
    var reader = new FileReader();
    reader.onload = function(e) {
      options.Custom_script_svg_value = e.target.result;
      console.log(options); // For demonstration, shows the updated options object
    };
      reader.readAsText(file); // For SVGs or other text-based files
    document.getElementById('Builtin_script_name').disabled = true;
    document.getElementById('customScriptLabel').style.fontWeight = 'bold';
    document.getElementById('customScriptLabel').style.color = '#007bff';
    document.getElementById('customScriptLabel').textContent = 'Custom script selected: '+  file.name;
      document.getElementById('Custom_script_svg_value').textContent = "change file"
  }});

// Add event listener to form elements to update options object
document.getElementById("form-container").addEventListener('change', updateOptions);

document.getElementById('renderButton').addEventListener('click', function() {
	Render(JSON.stringify(options))
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
  populateLanguages();
});
