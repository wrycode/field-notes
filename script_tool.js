var options = {
    "Debug": false,
    "Draw_bounding_boxes": false,
    "Form_stroke_width": 0.35,
    "Custom_script_svg_value":  "",    
    "Builtin_script_name": "teen_script.svg",
    "Language_code": "en_US",
    "Image_width": 500,
    "Input_text": "Let's start a\nnew paragraph. Here's a bunch of random text. How does it look?",
    "Space_between_metaforms": 5,
    "Space_between_lines": 10,
    "Margin": 5
};

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
    console.log(options);
}

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

document.getElementById('defaultsButton').addEventListener('click', function() {
    setDefaultValues()
});


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

// setDefaultValues();

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

// const logTextarea = document.getElementById('log');
// logTextarea.addEventListener('input', function() {
//     this.style.height = 'auto';
//     this.style.height = this.scrollHeight + 'px';
// });

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
