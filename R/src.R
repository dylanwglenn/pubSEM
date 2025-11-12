# Some useful keyboard shortcuts for package authoring:
#
#   Install Package:           'Ctrl + Shift + B'
#   Check Package:             'Ctrl + Shift + E'
#   Test Package:              'Ctrl + Shift + T'


#' Open the pubSEM GUI editor
#'
#' @param fit A lavaan object
#' @param standardized A bool value
#' @param project_name A string denoting the name of the pubSEM project to which
#'   a layout should be stored. A project file stores persistent data for
#'   reproducible diagrams — project names should be unique EVEN ACROSS R
#'   PROJECTS!
#' @returns nothing
#' @export
edit_path_diagram <- function(fit, standardized, project_name) {
    df_fit <- extract_lavaan_params(fit, standardized)

    usr_data_dir <- tools::R_user_dir("pkg", which = "data")
    base_dir <- paste0(base_dir, "/pubSEM/")
    file_path <- paste0(base_dir, "temp.json")

    jsonlite::write_json(df_fit,
        path = file_path,
        pretty = TRUE
    )

    if (Sys.info()['sysname'] == "Windows") {
        gui_exec_path <- system.file("main.exe", package = "pubSEM")
    } else {
        gui_exec_path <- system.file("main", package = "pubSEM")
    }

    # run the GUI executable
    system2(gui_exec_path,
            args = c(base_dir, project_name)
    )
}


# Get selected data from a lavaan fit
# returns a data frame
extract_lavaan_params <- function(fit, standardized) {
    df_paramTable <- lavaan::parametertable(fit)

    if (standardized) {
        df_estimates <- lavaan::standardizedsolution(fit)
        names(df_estimates)[names(df_estimates) == "est.std"] <- "est"
    } else {
        df_estimates <- lavaan::parameterestimates(fit)
    }

    # select desired columns from fit data
    df_paramTable_filtered <- df_paramTable[, c("lhs",
                                                "op",
                                                "rhs",
                                                "user",
                                                "label",
                                                "group")]

    df_estimates_filtered <- df_estimates[, c("lhs",
                                              "op",
                                              "rhs",
                                              "est",
                                              "pvalue",
                                              "ci.lower",
                                              "ci.upper")]

    # merge the two tables into one
    df <- merge(df_paramTable_filtered, df_estimates_filtered)
    # change problematic names
    names(df)[names(df) == "ci.lower"] <- "ci_lower"
    names(df)[names(df) == "ci.upper"] <- "ci_upper"
    return(df)
}
