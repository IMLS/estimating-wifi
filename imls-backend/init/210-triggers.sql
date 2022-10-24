-- migrate:up

--
-- Name: pgrst_watch(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.pgrst_watch() RETURNS event_trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NOTIFY pgrst, 'reload schema';
END;
$$;


--
-- Name: pgrst_watch; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER pgrst_watch ON ddl_command_end
   EXECUTE FUNCTION public.pgrst_watch();

-- migrate:down
