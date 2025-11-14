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
#'   reproducible diagrams in the current working directoloadry — project names
#'   should be unique EVEN ACROSS R PROJECTS!
#' @returns nothing
#' @export
sem_gui <- function(fit, project_name, standardized = FALSE) {
    df_fit <- extract_lavaan_params(fit, standardized)

    base_dir <- tools::R_user_dir("pubSEM", which = "data")
    file_path <- file.path(base_dir, "temp.json")

    if (!dir.exists(file.path(base_dir))) {
        dir.create(file.path(base_dir), recursive = TRUE)
    }

    jsonlite::write_json(df_fit,
        path = file_path,
        pretty = TRUE
    )

    if (Sys.info()['sysname'] == "Windows") {
        gui_exec_path <- system.file("bin", "sem_gui.exe", package = "pubSEM", mustWork = TRUE)
    } else {
        gui_exec_path <- system.file("bin", "sem_gui", package = "pubSEM", mustWork = TRUE)
    }

    # run the GUI executable
    system2(gui_exec_path,
            args = c(base_dir, project_name, "edit")
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
#
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
