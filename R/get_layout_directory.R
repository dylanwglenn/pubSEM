#' Get the directory holding the layout files
#'
#' @returns a file path to the directory containing pubSEM layout files
#' @export
get_layout_directory <- function() {
    return(tools::R_user_dir("pubSEM", which = "data"))
}
