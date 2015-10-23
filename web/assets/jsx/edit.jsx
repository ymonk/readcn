/**
 * Created by myan on 21/09/15.
 */

var apiHost = "http://writeuptube.com:5050";
var apiKeyPostfix = "key=abc123";

var getURLParameter = function (sParam) {
  var sPageURL = window.location.search.substring(1);
  var sURLVariables = sPageURL.split('&');
  for (var i = 0; i < sURLVariables.length; i++) {
    var sParameterName = sURLVariables[i].split('=');
    if (sParameterName[0] == sParam) {
      return sParameterName[1];
    }
  }
};

// Get permalink from URL like http://localhost:5000/article?v=55e9423f5485093cbdfe835f
var permalinkFromURL = function () {
  return getURLParameter("v")
};

var postTargetURL = function () {
  return apiHost + "/article/" + permalinkFromURL() + "?" + apiKeyPostfix;
};


var EditPane = React.createClass({
  getInitialState: function () {
    return {article: {}};
  },

  componentDidMount: function () {
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      cache: false,
      success: function (data) {
        this.setState({article: data})
      }.bind(this),
      error: function (xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },

  onReset: function (evt) {
    var fullURL = apiHost + "/article/" + this.currentPermalink() + "?" + apiKeyPostfix;
    $.ajax({
      url: fullURL,
      dataType: 'json',
      cache: false,
      success: function (data) {
        this.setState({article: data})
      }.bind(this),
      error: function (xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
    evt.preventDefault();
  },

  currentPermalink: function () {
    if (typeof(this.state.article.permalink) == "undefined") {
      return permalinkFromURL();
    }
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

  onAuthorChange: function (evt) {
    this.handleChange("author", evt.target.value);
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

  onBodyChange: function (evt) {
    this.handleChange("body", evt.target.value);
  },

  onPermalinkChange: function (evt) {
    this.handleChange("permalink", evt.target.value);
  },

  onSave: function (evt) {
    var url = postTargetURL();
    $.post(url, JSON.stringify(this.state.article)).done(function (d, s, r) {
      var newURL =  "/edit?v=" + this.currentPermalink();
      history.pushState({}, "", newURL);
      $("#show-message").trigger("click");
    }.bind(this));
    evt.preventDefault();
  },

  onSaveClose: function (evt) {
    var url = postTargetURL();
    $.post(url, JSON.stringify(this.state.article)).done(function (d, s, r) {
      history.back();
    }.bind(this));
    evt.preventDefault();
  },

  onDelete: function (evt) {
    var url = apiHost + "/delete/article/" + permalinkFromURL() + "?" + apiKeyPostfix;
    $.ajax({
      url: url,
      dataType: 'json',
      cache: false,
      success: function () {
        var newURL =  "/edit";
        window.location = newURL;
      }.bind(this),
      error: function (xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
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
              <span className="caption-subject font-green-sharp bold uppercase"> Edit Article </span>
            </div>
            <div className="btn-set pull-right">
              <button className="btn btn-default btn-circle" onClick={this.onReset}><i className="fa fa-reply"></i> Reset</button>
              <button className="btn green-haze btn-circle" onClick={this.onSaveClose}><i className="fa fa-check-circle"></i> Save and Back</button>
              <button className="btn green-haze btn-circle" onClick={this.onSave}><i className="fa fa-check-circle"></i> Save and Continue</button>
              <button className="btn red-haze btn-circle" onClick={this.onDelete}><i className="fa fa-times-circle"></i> Delete</button>
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
              <label className="col-md-2 control-label">Author: <span className="required"> * </span></label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[author]" onChange={this.onAuthorChange}
                       placeholder="author" value={this.state.article.author} />
              </div>
            </div>

            <div className="form-group">
              <label className="col-md-2 control-label">
                Summary: <span className="required"> * </span>
              </label>

              <div className="col-md-10">
                <textarea className="form-control" name="article[preview]" onChange={this.onPreviewChange}
                          style={previewStyle} value={this.state.article.preview} >
                </textarea>
              </div>
            </div>

            <div className="form-group">
              <label className="col-md-2 control-label">HSK Character Level: <span className="required">
                                                                              * </span>
              </label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[char_level]" onChange={this.onCharLevelChange}
                       placeholder="" value={this.state.article.char_level} />
              </div>
            </div>

            <div className="form-group">
              <label className="col-md-2 control-label">HSK Vocabulary Level: <span className="required">
                  * </span>
              </label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[vocabulary_level]" onChange={this.onVocabularyLevelChange}
                       placeholder="" value={this.state.article.vocabulary_level} />
              </div>
            </div>
            <div className="form-group">
              <label className="col-md-2 control-label">HSK Grammar Level: <span className="required">
                                                                              * </span>
              </label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[grammar_level]" onChange={this.onGrammarLevelChange}
                       placeholder="" value={this.state.article.grammar_level} />
              </div>
            </div>

            <div className="form-group">
              <label className="col-md-2 control-label">Source: <span className="required">
                                                                              * </span>
              </label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[source]" onChange={this.onSourceChange}
                       placeholder="" value={this.state.article.source}/>
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
              <label className="col-md-2 control-label">Text: <span className="required">
                                                                              * </span>
              </label>

              <div className="col-md-10">
                <textarea className="form-control" style={bodyStyle} name="article[body]" onChange={this.onBodyChange}
                          value={rawMarkup} />
              </div>
            </div>


            <div className="form-group">
              <label className="col-md-2 control-label">Permanent Link: <span className="required"> * </span>
              </label>

              <div className="col-md-10">
                <input type="text" className="form-control" name="article[permalink]" onChange={this.onPermalinkChange}
                       placeholder="" value={this.state.article.permalink} />
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
    <EditPane url={postTargetURL()} />,
    document.getElementById('content')
  );
});


