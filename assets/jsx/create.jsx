/**
 * Created by myan on 21/09/15.
 */

var apiHost = "http://writeuptube.com:5050";
var apiKeyPostfix = "key=abc123";

//var getURLParameter = function (sParam) {
//  var sPageURL = window.location.search.substring(1);
//  var sURLVariables = sPageURL.split('&');
//  for (var i = 0; i < sURLVariables.length; i++) {
//    var sParameterName = sURLVariables[i].split('=');
//    if (sParameterName[0] == sParam) {
//      return sParameterName[1];
//    }
//  }
//};

var actionTargetURL = function () {
  return apiHost + "/create/article" + "?" + apiKeyPostfix;
};

var createIdURL = function () {
  return apiHost + "/create/id" + "?" + apiKeyPostfix;
};

var EditPane = React.createClass({
  getInitialState: function () {
    return {article: {
      visibility: 2
    }};
  },

  // Get a new ObjectId as the auto-generated
  componentDidMount: function () {
    var url = createIdURL();
    $.ajax({
      url: url,
      dataType: 'json',
      cache: false,
      success: function (data) {
        var article = {
          _id: data.id,
          permalink: data.permalink,
          visibility: 2
        };
        this.setState({article: article})
      }.bind(this),
      error: function (xhr, status, err) {
        console.error(url, status, err.toString());
      }.bind(this)
    });
  },

  currentPermalink: function() {
    return this.state.article.permalink;
  },

  handleChange: function (key, val) {
    var article = this.state.article;
    article[key] = val;
    this.setState({article: article});
  },

  onTitleChange: function (evt) {
    this.handleChange("title", evt.target.value);
  },

  onPreviewChange: function (evt) {
    this.handleChange("preview", evt.target.value);
  },

  onCharLevelChange: function (evt) {
    this.handleChange("char_level", evt.target.value);
  },

  onVocabularyLevelChange: function (evt) {
    this.handleChange("vocabulary_level", evt.target.value);
  },

  onGrammarLevelChange: function (evt) {
    this.handleChange("grammar_level", evt.target.value);
  },

  onSourceChange: function (evt) {
    this.handleChange("source", evt.target.value);
  },

  onCategoriesChange: function (evt) {
    this.handleChange("categories", evt.target.value);
  },

  onVisibilityChange: function (evt) {
    var visibility = 2;
    if (evt.target.checked == false) {
      visibility = 1;
    }
    var article = this.state.article;
    article.visibility = visibility;
    this.setState({article: article});
  },

  onBodyChange: function (evt) {
    this.handleChange("body", evt.target.value);
  },

  onPermalinkChange: function (evt) {
    this.handleChange("permalink", evt.target.value);
  },

  onSaveClose: function (evt) {
    var url = actionTargetURL();
    console.log("About to POST to", url);
    $.post(url, JSON.stringify(this.state.article)).done(function (d, s, r) {
      console.log("Jump to", "/static/read.html?v=" + this.currentPermalink());
      window.location = "/static/read.html?v=" + this.currentPermalink();
    }.bind(this));
    evt.preventDefault();
  },

  onCancel: function (evt) {
    this.setState({article: {
      visibility: 2,
      preview: "",
      body: ""
    }});
    evt.preventDefault();
  },

  render: function () {
    var bodyStyle = {
      height: "400px"
    };

    var previewStyle = {
      height: "140px"
    };

    var rawMarkup = this.state.article.body;
    return (
      <form className="form-horizontal form-row-seperated">
        <div className="portlet light">
          <div className="portlet-title">
            <div className="caption">
              <i className="icon-edit font-green-sharp"></i>
              <span className="caption-subject font-green-sharp bold uppercase"> Create Article </span>
            </div>
            <div className="btn-set pull-right">
              <button className="btn green-haze btn-circle" onClick={this.onSaveClose}><i className="fa fa-check-circle"></i> Save and View</button>
              <button className="btn red-haze btn-circle" onClick={this.onCancel}><i className="fa fa-times-circle"></i> Cancel</button>
            </div>
          </div>
          <div className="form-body">
            <div className="form-group">
              <label className="col-md-2 control-label">Title: <span className="required"> * </span></label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[title]" onChange={this.onTitleChange}
                       placeholder="title" value={this.state.article.title} />
              </div>
            </div>
            <div className="form-group">
              <label className="col-md-2 control-label">
                Summary: <span className="required"> * </span>
              </label>

              <div className="col-md-10">
                <textarea className="form-control" name="article[preview]" onChange={this.onPreviewChange}
                          style={previewStyle} placeholder="summary, less than 100 characters"
                          value={this.state.article.preview} >
                </textarea>
              </div>
            </div>

            <div className="form-group">
              <label className="col-md-2 control-label">HSK Character Level: <span className="required">
                                                                              * </span>
              </label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[char_level]" onChange={this.onCharLevelChange}
                       placeholder="1" value={this.state.article.char_level} />
              </div>
            </div>

            <div className="form-group">
              <label className="col-md-2 control-label">HSK Vocabulary Level: <span className="required">
                  * </span>
              </label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[vocabulary_level]" onChange={this.onVocabularyLevelChange}
                       placeholder="1" value={this.state.article.vocabulary_level} />
              </div>
            </div>
            <div className="form-group">
              <label className="col-md-2 control-label">HSK Grammar Level: <span className="required">
                                                                              * </span>
              </label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[grammar_level]" onChange={this.onGrammarLevelChange}
                       placeholder="1" value={this.state.article.grammar_level} />
              </div>
            </div>

            <div className="form-group">
              <label className="col-md-2 control-label">Source: <span className="required">
                                                                              * </span>
              </label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[source]" onChange={this.onSourceChange}
                       placeholder="source" value={this.state.article.source}/>
              </div>
            </div>

            <div className="form-group">
              <label className="col-md-2 control-label">Categories: <span className="required">
                                                                              * </span>
              </label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[categories]" onChange={this.onCategoriesChange}
                       placeholder="Categories, seperated with comma or semicomma" value={this.state.article.categories}/>
              </div>
            </div>

            <div className="form-group">
              <label className="col-md-2 control-label">Public: <span className="required">
                                                                              * </span>
              </label>

              <div className="col-md-1">
                <input type="checkbox" className="form-control" name="article[visibility]" onChange={this.onVisibilityChange}
                       checked={this.state.article.visibility == 2} id="public-check" />
              </div>
            </div>

            <div className="form-group">
              <label className="col-md-2 control-label">Text: <span className="required">
                                                                              * </span>
              </label>

              <div className="col-md-10">
                <textarea className="form-control" style={bodyStyle} name="article[body]" onChange={this.onBodyChange}
                          placeholder="content" value={rawMarkup} />
              </div>
            </div>


            <div className="form-group">
              <label className="col-md-2 control-label">Permanent Link:
              </label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[permalink]" onChange={this.onPermalinkChange}
                       placeholder="Auto generated" value={this.state.article.permalink} />
              </div>
            </div>
          </div>
        </div>
      </form>
    );
  }
});

$(function () {
  React.render(
    <EditPane url={actionTargetURL()} />,
    document.getElementById('content')
  );
});


